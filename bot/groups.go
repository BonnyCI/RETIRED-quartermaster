package bot

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/lib"
	"github.com/bonnyci/quartermaster/web/client"
)

func GroupsHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In GroupsHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "groups",
			Usage: "groups",
			Short: "List groups.",
		},
		HelpFmt{
			Cmd:   "groups add",
			Usage: "groups add <group(s)>",
			Short: "Add a new group(s). (groups as a comma separated list)",
		},
		HelpFmt{
			Cmd:   "groups del",
			Usage: "groups del <group(s)>",
			Short: "Delete a group(s). (groups as a comma separated list)",
		},
		HelpFmt{
			Cmd:   "groups addmembers",
			Usage: "groups addmembers <group(s)> <user(s)>",
			Short: "Add user(s) to a group(s) member list. (groups and users as comma separated list)",
		},
		HelpFmt{
			Cmd:   "groups delmembers",
			Usage: "groups delmembers <group(s)> <user(s)>",
			Short: "Del users(s) from a group(s) member list. (groups and users as comma separated list)",
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func GroupsBase(i *Irc, command *Command) {
	jww.DEBUG.Println("Listing Groups")
	gs, err := client.GetGroups()
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	for _, g := range gs {
		i.Sendf(command.Target, "Group: %s", g.String())
	}
}

func GroupsAdd(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		GroupsHelp(i, command)
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
	jww.DEBUG.Printf("Adding Group(s): %+v", command.Args)

	for _, v := range strings.Split(command.Args[0], ",") {
		g, err := client.AddGroup(i.Api.User, i.Api.Pass, v)
		if err != nil {
			efmt := "Failed to create group: %s"
			jww.ERROR.Printf(efmt, v)
			i.Sendf(command.Target, efmt, v)
		}
		jww.DEBUG.Printf("Group: %+v", g)
		i.Sendf(command.Target, "Group %s created.", v)
	}
}

func GroupsDel(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		GroupsHelp(i, command)
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
	jww.DEBUG.Printf("Deleting Group(s): %+v", command.Args)
	groups := strings.Split(command.Args[0], ",")
	if fnd, _ := lib.In("Admin", groups); fnd {
		efmt := "User (%s) is not authorized to delete Admin group."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	for _, v := range groups {
		err := client.DelGroup(i.Api.User, i.Api.Pass, v)
		if err != nil {
			efmt := "Failed to create group: %s"
			jww.ERROR.Printf(efmt, v)
			i.Sendf(command.Target, efmt, v)
		}
	}
}

func GroupsAddMembers(i *Irc, command *Command) {
	jww.DEBUG.Printf("GAM: %+v", command)
	jww.DEBUG.Printf("Args: %d", len(command.Args))
	if len(command.Args) != 2 {
		GroupsHelp(i, command)
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
	jww.DEBUG.Printf("Adding Members(s) to Group(s): %+v to %+v", command.Args[1], command.Args[0])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")

	for _, group := range gs {
		for _, user := range us {
			_, err := client.AddMemberToGroup(i.Api.User, i.Api.Pass, group, user)
			if err != nil {
				efmt := "Could not add %s to %s."
				jww.ERROR.Printf(efmt, user, group)
				i.Sendf(command.Target, efmt, user, group)
				return
			}
		}
	}
}

func GroupsDelMembers(i *Irc, command *Command) {
	if len(command.Args) != 2 {
		GroupsHelp(i, command)
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
	jww.DEBUG.Printf("Deleting User(s) from Group(s): %+v to %+v", command.Args[0], command.Args[1])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")

	for _, group := range gs {
		for _, user := range us {
			err := client.DelMemberFromGroup(i.Api.User, i.Api.Pass, group, user)
			if err != nil {
				efmt := "Could not remove %s from %s."
				jww.ERROR.Printf(efmt, user, group)
				i.Sendf(command.Target, efmt, user, group)
				return
			}
		}
	}
}

func Groups(i *Irc, command *Command) {
	jww.DEBUG.Println("In Users")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("groups", GroupsBase)
	c.HandleFunc("add", GroupsAdd)
	c.HandleFunc("del", GroupsDel)
	c.HandleFunc("addmembers", GroupsAddMembers)
	c.HandleFunc("delmembers", GroupsDelMembers)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
