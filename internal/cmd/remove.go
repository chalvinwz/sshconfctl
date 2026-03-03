package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chalvinwz/sshconfctl/internal/prompt"
	"github.com/chalvinwz/sshconfctl/internal/sshconfig"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <host>",
	Short: "Remove a host entry (asks whether to backup first)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]

		cfg, err := sshconfig.Load()
		if err != nil {
			return err
		}

		if !sshconfig.HasHost(cfg, host) {
			return fmt.Errorf("host %q not found", host)
		}

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("\nYou are about to remove host %q from ~/.ssh/config\n", host)

		doBackup, err := prompt.AskYesNo(scanner, "Create backup first?", true)
		if err != nil {
			return err
		}
		if doBackup {
			if err := sshconfig.Backup(); err != nil {
				fmt.Printf("Warning: Backup failed: %v\n", err)
			} else {
				fmt.Println("Backup created → ~/.ssh/config.bak.YYYYMMDD_HHMMSS")
			}
		}

		fmt.Printf("Type the host alias to confirm deletion (%s): ", host)
		confirmText := strings.TrimSpace(prompt.Ask(scanner, ""))
		if confirmText != host {
			fmt.Println("Removal cancelled (confirmation mismatch).")
			return nil
		}

		if err := sshconfig.RemoveHost(cfg, host); err != nil {
			return err
		}
		if err := sshconfig.Save(cfg); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("→ Removed host %q successfully\n", host)
		return nil
	},
}
