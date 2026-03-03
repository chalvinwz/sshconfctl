package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chalvinwz/sshconfctl/internal/config"
	"github.com/chalvinwz/sshconfctl/internal/prompt"
	"github.com/chalvinwz/sshconfctl/internal/sshconfig"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <host>",
	Short: "Edit an existing host entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		def := config.LoadDefaults()
		scanner := bufio.NewScanner(os.Stdin)

		cfg, err := sshconfig.Load()
		if err != nil {
			return err
		}

		if !sshconfig.HasHost(cfg, host) {
			return fmt.Errorf("host %q not found", host)
		}

		currentHostname, _ := cfg.Get(host, "HostName")
		currentPort, _ := cfg.Get(host, "Port")
		if strings.TrimSpace(currentPort) == "" {
			currentPort = "22"
		}
		currentUser, _ := cfg.Get(host, "User")
		currentIdfile, _ := cfg.Get(host, "IdentityFile")

		fmt.Printf("\nEditing host: %s\n", host)

		hostname := prompt.Ask(scanner, fmt.Sprintf("Hostname / IP  [%s]: ", currentHostname))
		hostname = strings.TrimSpace(hostname)
		if hostname == "" {
			hostname = currentHostname
		}
		if hostname == "" {
			return fmt.Errorf("hostname is required")
		}
		if err := sshconfig.ValidateHostNameOrIP(hostname); err != nil {
			return err
		}

		port := prompt.Ask(scanner, fmt.Sprintf("Port           [%s]: ", currentPort))
		port = strings.TrimSpace(port)
		if port == "" {
			port = currentPort
		}
		if err := sshconfig.ValidatePort(port); err != nil {
			return err
		}

		user := prompt.Ask(scanner, fmt.Sprintf("User           [%s]: ", currentUser))
		user = strings.TrimSpace(user)
		if user == "" {
			user = currentUser
		}
		if user == "" {
			user = def.User
		}

		idfile := prompt.Ask(scanner, fmt.Sprintf("IdentityFile   [%s]: ", currentIdfile))
		idfile = strings.TrimSpace(idfile)
		if idfile == "" {
			idfile = currentIdfile
		}
		if idfile == "" {
			idfile = def.IdentityFile
		}

		if err := sshconfig.UpdateHost(cfg, host, hostname, port, user, idfile); err != nil {
			return err
		}
		if err := sshconfig.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("→ Updated host %q successfully\n", host)
		return nil
	},
}
