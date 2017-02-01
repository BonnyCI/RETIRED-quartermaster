package bot

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

func Quit(i *Irc, command *Command) {
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered to perform actions on this bot"
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	if fnd := lib.UserInGroup("Admin", u); fnd == false {
		efmt := "User (%s) is not authorized to perform this action."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	jww.DEBUG.Println("Quiting IRC")
	i.Disconnect()
}
