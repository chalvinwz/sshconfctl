package cmd

import (
	"fmt"

	"github.com/chalvinwz/sshconfctl/internal/sshconfig"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a timestamped backup of ~/.ssh/config",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sshconfig.Backup(); err != nil {
			return err
		}
		fmt.Println("Backup created successfully")
		return nil
	},
}
