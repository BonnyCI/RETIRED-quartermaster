package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/bonnyci/quartermaster/lib"
)

type Group struct {
	Name  string
	Users []string
}

type Data struct {
	Users  []string
	Groups []Group
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := InitializeConfig()
		if err != nil {
			return err
		}

		var dat Data
		if err := viper.UnmarshalKey("Data", &dat); err != nil {
			panic(err)
		}

		lib.AddUsers(dat.Users)

		for _, v := range dat.Groups {
			lib.AddGroups([]string{v.Name})
			lib.AddUsersToGroups([]string{v.Name}, v.Users)
		}

		for _, v := range dat.Groups {
			g, err := lib.GetGroup(v.Name)
			if err != nil {
				jww.ERROR.Println(err)
			}
			jww.DEBUG.Printf("g: %s", g.String())
		}

		return nil
	},
}

func init() {
}
