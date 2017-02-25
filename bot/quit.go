package bot

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/client"
)

func Quit(i *Irc, command *Command) {
	u, err := database.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered to perform actions on this bot"
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	g, _ := client.GetGroup("Admin")
	if fnd := database.UserInGroup(g, u); fnd == false {
		efmt := "User (%s) is not authorized to perform this action."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	jww.DEBUG.Println("Quiting IRC")
	i.Disconnect()
}
