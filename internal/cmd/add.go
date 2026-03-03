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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new host entries interactively (type 'exit' to quit)",
	RunE: func(cmd *cobra.Command, args []string) error {
		def := config.LoadDefaults()
		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Println("\n─ Add new host ──────────────────────────────────────── (type 'exit' to quit)")

			host := prompt.Ask(scanner, "Host alias     : ")
			host = strings.TrimSpace(host)
			if host == "exit" {
				break
			}
			if host == "" {
				continue
			}

			cfg, err := sshconfig.Load()
			if err != nil {
				return err
			}

			if sshconfig.HasHost(cfg, host) {
				fmt.Println("Error: Host alias already exists")
				continue
			}

			if err := sshconfig.ValidateAlias(host); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			hostname := prompt.Ask(scanner, "Hostname / IP  : ")
			hostname = strings.TrimSpace(hostname)
			if hostname == "" {
				fmt.Println("Error: Hostname / IP is required")
				continue
			}

			if err := sshconfig.ValidateHostNameOrIP(hostname); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			port := prompt.Ask(scanner, "Port           [22]: ")
			port = strings.TrimSpace(port)
			if port == "" {
				port = "22"
			}
			if err := sshconfig.ValidatePort(port); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			user := prompt.Ask(scanner, fmt.Sprintf("User           [%s]: ", def.User))
			user = strings.TrimSpace(user)
			if user == "" {
				user = def.User
			}

			idfile := prompt.Ask(scanner, fmt.Sprintf("IdentityFile   [%s]: ", def.IdentityFile))
			idfile = strings.TrimSpace(idfile)
			if idfile == "" {
				idfile = def.IdentityFile
			}

			if err := sshconfig.AppendHost(cfg, host, hostname, port, user, idfile); err != nil {
				return err
			}
			if err := sshconfig.Save(cfg); err != nil {
				return err
			}

			fmt.Printf("\n→ Added host '%s' successfully\n", host)
		}

		return nil
	},
}
