package cmd

import (
	"github.com/spf13/cobra"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/web/client"
	"github.com/bonnyci/quartermaster/web/endpoints"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := InitializeConfig()
		if err != nil {
			return err
		}

		endpoints.BackupHTTP()
		endpoints.ApiHTTP()

		us, _ := client.GetUsers()
		jww.WARN.Printf("users: %+v", us)
		u, err := client.AddUser("quartermaster", "quartermaster", "test")
		jww.FATAL.Println(err)
		jww.WARN.Printf("U1: %+v", u)
		u, err = client.GetUser("test")
		jww.FATAL.Println(err)
		jww.WARN.Printf("U2: %+v", u)
		client.DelUser("quartermaster", "quartermaster", "test")

		for {
		}

		return nil
	},
}

func init() {
}
