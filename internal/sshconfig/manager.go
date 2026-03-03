package sshconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kevinburke/ssh_config"
)

const ConfigPath = "~/.ssh/config"

func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[2:])
	}
	return p
}

func Load() (*ssh_config.Config, error) {
	path := ExpandPath(ConfigPath)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			return nil, err
		}
		f, err = os.Open(path)
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ssh_config.Decode(f)
}

// Save writes ~/.ssh/config with consistent two-space indentation.
func Save(cfg *ssh_config.Config) error {
	path := ExpandPath(ConfigPath)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	content := normalizeIndentation(cfg.String())

	if _, err := io.WriteString(f, content); err != nil {
		return err
	}
	return nil
}

func HasHost(cfg *ssh_config.Config, alias string) bool {
	aliasLower := strings.ToLower(alias) // optional: warn on case issues later
	for _, h := range cfg.Hosts {
		for _, pat := range h.Patterns {
			patStr := pat.String()
			if patStr == alias || strings.ToLower(patStr) == aliasLower {
				if patStr != alias {
					fmt.Printf("Warning: Found similar alias with different case: %q\n", patStr)
				}
				return true
			}
		}
	}
	return false
}

func GetAllHosts(cfg *ssh_config.Config) []string {
	var hosts []string
	for _, h := range cfg.Hosts {
		for _, pat := range h.Patterns {
			s := pat.String()
			if s != "*" && s != "" {
				hosts = append(hosts, s)
			}
		}
	}
	return hosts
}

func AppendHost(cfg *ssh_config.Config, host, hostname, port, user, idfile string) error {
	pat, err := ssh_config.NewPattern(host)
	if err != nil {
		return fmt.Errorf("invalid host pattern %q: %w", host, err)
	}

	h := &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{pat},
		Nodes: []ssh_config.Node{
			&ssh_config.KV{Key: "HostName", Value: hostname},
			&ssh_config.KV{Key: "Port", Value: port},
			&ssh_config.KV{Key: "User", Value: user},
			&ssh_config.KV{Key: "PreferredAuthentications", Value: "publickey"},
			&ssh_config.KV{Key: "IdentityFile", Value: idfile},
		},
	}
	cfg.Hosts = append(cfg.Hosts, h)
	return nil
}

func UpdateHost(cfg *ssh_config.Config, host, hostname, port, user, idfile string) error {
	for _, h := range cfg.Hosts {
		for _, pat := range h.Patterns {
			if pat.String() == host {
				setKV(h, "HostName", hostname)
				setKV(h, "Port", port)
				setKV(h, "User", user)
				setKV(h, "PreferredAuthentications", "publickey")
				setKV(h, "IdentityFile", idfile)
				return nil
			}
		}
	}
	return fmt.Errorf("host %q not found", host)
}

func setKV(h *ssh_config.Host, key, value string) {
	for i, n := range h.Nodes {
		if kv, ok := n.(*ssh_config.KV); ok && kv.Key == key {
			h.Nodes[i] = &ssh_config.KV{Key: key, Value: value}
			return
		}
	}
	h.Nodes = append(h.Nodes, &ssh_config.KV{Key: key, Value: value})
}

func RemoveHost(cfg *ssh_config.Config, name string) error {
	for i, h := range cfg.Hosts {
		for _, p := range h.Patterns {
			if p.String() == name {
				cfg.Hosts = append(cfg.Hosts[:i], cfg.Hosts[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("host %q not found", name)
}

func normalizeIndentation(content string) string {
	content = strings.ReplaceAll(content, "\t", "  ")
	lines := strings.Split(content, "\n")

	out := make([]string, 0, len(lines))
	seenHost := false

	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "#") {
			out = append(out, trim)
			continue
		}
		if strings.HasPrefix(trim, "Host ") || strings.HasPrefix(trim, "Match ") {
			if seenHost && len(out) > 0 && out[len(out)-1] != "" {
				out = append(out, "")
			}
			out = append(out, trim)
			seenHost = true
			continue
		}
		out = append(out, "  "+trim)
	}

	return strings.TrimRight(strings.Join(out, "\n"), "\n\t \r") + "\n"
}

func Backup() error {
	src := ExpandPath(ConfigPath)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // nothing to backup
	}

	bak := fmt.Sprintf("%s.bak.%s", src, time.Now().Format("20060102_150405"))
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(bak, data, 0600)
}
