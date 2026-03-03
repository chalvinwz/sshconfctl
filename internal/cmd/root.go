package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chalvinwz/sshconfctl/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "sshconfctl",
	Short: "Manage SSH config entries easily",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/sshconfctl/config.yaml)")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := filepath.Join(home, ".config", "sshconfctl")
		cobra.CheckErr(os.MkdirAll(cfgDir, 0o700))
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			defs := config.LoadDefaults()
			if isInteractiveTerminal() {
				fmt.Println("First run setup: config.yaml not found.")
				defs.User = askRequired("Default user")
				defs.IdentityFile = askRequired("Default identity file")
			}
			cobra.CheckErr(config.SaveDefaults(defs))
			return
		}
		cobra.CheckErr(err)
	}
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func askRequired(label string) string {
	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", label)
		if !s.Scan() {
			return ""
		}
		v := strings.TrimSpace(s.Text())
		if v != "" {
			return v
		}
		fmt.Println("Value cannot be empty.")
	}
}
