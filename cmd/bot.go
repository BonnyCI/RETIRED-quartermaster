package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/bonnyci/quartermaster/bot"
	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/endpoints"
)

// getCmd represents the get command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the bot",
	Long:  `Start the bot.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := InitializeConfig()
		if err != nil {
			return err
		}

		database.GetInstance()

		jww.INFO.Println("GET - Server: ", viper.GetString("server"))
		jww.INFO.Println("GET - Port: ", viper.GetString("port"))
		jww.INFO.Println("GET - Debug: ", viper.GetString("debug"))

		endpoints.BackupHTTP()
		endpoints.ApiHTTP()

		i := bot.GetIrc()
		bot.Configure(i)
		i.Connect()
		defer database.CloseInstance()

		return nil
	},
}

func init() {
}
