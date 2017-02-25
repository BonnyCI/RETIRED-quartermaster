package bot

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/client"
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
	us, err := client.GetUsers()
	if err != nil {
		jww.ERROR.Println(err)
		i.Sendf(command.Target, "Command error: %s", err.Error())
	}

	for _, u := range us {
		i.Sendf(command.Target, "User: %s", u.String())
	}
}

func UsersAdd(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		UsersHelp(i, command)
		return
	}

	u, err := client.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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

	jww.DEBUG.Printf("Adding User(s): %+v", command.Args)
	for _, v := range strings.Split(command.Args[0], ",") {
		_, err := client.AddUser(i.Api.User, i.Api.Pass, v)
		if err != nil {
			efmt := "Failed to create user: %s"
			jww.ERROR.Printf(efmt, v)
			i.Sendf(command.Target, efmt, v)
		}
	}
}

func UsersDel(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		UsersHelp(i, command)
		return
	}
	u, err := client.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Deleting User(s): %+v", command.Args)

	for _, v := range strings.Split(command.Args[0], ",") {
		err := client.DelUser(i.Api.User, i.Api.Pass, v)
		if err != nil {
			efmt := "Failed to delete user: %s"
			jww.ERROR.Printf(efmt, v)
			i.Sendf(command.Target, efmt, v)
		}
	}
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
