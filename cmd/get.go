package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
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

		lib.AddGroups([]string{"Admin"})
		lib.AddUsers([]string{"phschwartz"})
		lib.AddUsersToGroups("Admin", []string{"Admin"}, []string{"phschwartz"})

		u, _ := lib.GetUser("phschwartz")
		st := lib.GetStatus(u, lib.DStamp)

		jww.DEBUG.Printf("u: %+v", u)
		jww.DEBUG.Printf("st: %+v", st)

		return nil
	},
}

func init() {
}
