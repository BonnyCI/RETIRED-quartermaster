package lib

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadGlobalConfig loads Hugo configuration into the global Viper.
func LoadGlobalConfig(configFilename string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("quartermaster")
	viper.SetConfigFile(configFilename)
	viper.AddConfigPath("/etc/quartermaster")
	viper.AddConfigPath("$HOME/.quartermaster")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		return fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `quartermaster help new` for details. (%s)\n", err)
	}

	loadDefaultSettings()

	return nil
}

func loadDefaultSettings() {
	viper.SetDefault("server", "irc.freenode.net")
	viper.SetDefault("port", 6697)
	viper.SetDefault("UseTLS", true)
	viper.SetDefault("user", "quatermaster")
	viper.SetDefault("nick", "quartermaster")
	viper.SetDefault("pass", "")
	viper.SetDefault("debug", false)
	viper.SetDefault("standupLogs", "logs/standup-{{date}}")
}
