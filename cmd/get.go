package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/database"
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

		db := &database.Database{Db: "test.db"}
		jww.DEBUG.Println("DB: ", db.Db)

		db.Open()
		defer db.Close()

		ps := database.UserS{Irc: "phschwartz", Gh: "pschwartz"}
		mm := database.UserS{Irc: "mms", Gh: "mms"}

		ps.Save(db)
		mm.Save(db)

		var qr database.UserS
		qr.GetOne(db, "Irc", "phschwartz")

		jww.DEBUG.Printf("%+v", qr)

		qr.GetOne(db, "Irc", "mms")
		jww.DEBUG.Printf("%+v", qr)

		var qa []database.UserS
		database.GetAll(db, &qa)
		jww.DEBUG.Printf("%+v", qa)

		return nil
	},
}

func init() {
}
