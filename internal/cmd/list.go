package cmd

import (
	"fmt"

	"github.com/chalvinwz/sshconfctl/internal/sshconfig"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all host aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := sshconfig.Load()
		if err != nil {
			return err
		}

		hosts := sshconfig.GetAllHosts(cfg)
		if len(hosts) == 0 {
			fmt.Println("No hosts defined in ~/.ssh/config")
			return nil
		}

		fmt.Println("Defined hosts:")
		for _, h := range hosts {
			fmt.Printf("  • %s\n", h)
		}

		return nil
	},
}
