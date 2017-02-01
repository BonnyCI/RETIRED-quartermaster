package bot

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

func UsersHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In UsersHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "users",
			Usage: "users",
			Short: "List all users",
		},
		HelpFmt{
			Cmd:   "users add",
			Usage: "users add <user|users>",
			Short: "Add a user or list of users (comma delimiated) to the bot.",
		},
		HelpFmt{
			Cmd:   "users del",
			Usage: "users del <user|users>",
			Short: "Delete a user or list of users (comma seperated) from the bot.",
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func UsersBase(i *Irc, command *Command) {
	jww.DEBUG.Println("Listing Users")
	us := lib.ListUsers()
	for _, u := range us {
		i.Sendf(command.Target, "User: %s", u.String())
	}
}

func UsersAdd(i *Irc, command *Command) {
	jww.DEBUG.Printf("Adding User(s): %+v", command.Args)
	lib.AddUsers(strings.Split(command.Args[0], ","))
}

func UsersDel(i *Irc, command *Command) {
	jww.DEBUG.Printf("Deleting User(s): %+v", command.Args)
	lib.DelUsers(strings.Split(command.Args[0], ","))
}

func Users(i *Irc, command *Command) {
	jww.DEBUG.Println("In Users")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("users", UsersBase)
	c.HandleFunc("add", UsersAdd)
	c.HandleFunc("del", UsersDel)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
