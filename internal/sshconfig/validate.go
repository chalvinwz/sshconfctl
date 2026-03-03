package sshconfig

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	hostAliasRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	fqdnRe      = regexp.MustCompile(`^(?i)([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)$`)
)

func ValidateAlias(alias string) error {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return fmt.Errorf("host alias is required")
	}
	if !hostAliasRe.MatchString(alias) {
		return fmt.Errorf("host alias can only contain letters, numbers, dots, dashes and underscores")
	}
	return nil
}

func ValidateHostNameOrIP(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return fmt.Errorf("hostname / IP is required")
	}
	if net.ParseIP(v) != nil {
		return nil
	}
	if fqdnRe.MatchString(v) && len(v) <= 253 && hasAlpha(v) {
		return nil
	}
	return fmt.Errorf("invalid format: only valid IP addresses or FQDNs are allowed")
}

func ValidatePort(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return fmt.Errorf("port is required")
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > 65535 {
		return fmt.Errorf("invalid port: must be number 1-65535")
	}
	return nil
}

func hasAlpha(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}
