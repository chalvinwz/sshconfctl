package cmd

import (
	"fmt"

	"github.com/chalvinwz/sshconfctl/internal/sshconfig"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Reformat ~/.ssh/config with consistent 2-space indentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := sshconfig.Load()
		if err != nil {
			return err
		}
		if err := sshconfig.Save(cfg); err != nil {
			return err
		}
		fmt.Println("SSH config formatted successfully.")
		return nil
	},
}
