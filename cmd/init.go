package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/bonnyci/quartermaster/database"
)

type Group struct {
	Name  string
	Users []string
}

type User struct {
	Username string
	Password string
}
type Data struct {
	Users  []User
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

		for _, v := range dat.Users {
			database.AddUsers([]string{v.Username})
			database.AddPassword(v.Username, v.Password)
		}

		for _, v := range dat.Groups {
			database.AddGroups([]string{v.Name})
			database.AddUsersToGroups([]string{v.Name}, v.Users)
		}

		for _, v := range dat.Groups {
			g, err := database.GetGroup(v.Name)
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
