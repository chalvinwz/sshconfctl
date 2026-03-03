package config

import "github.com/spf13/viper"

type Defaults struct {
	User         string
	IdentityFile string
}

func LoadDefaults() Defaults {
	return Defaults{
		User:         viper.GetString("defaults.user"),
		IdentityFile: viper.GetString("defaults.identity_file"),
	}
}

func SaveDefaults(d Defaults) error {
	viper.Set("defaults.user", d.User)
	viper.Set("defaults.identity_file", d.IdentityFile)
	if err := viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfig()
		}
		return err
	}
	return nil
}
