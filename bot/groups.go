package bot

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

func GroupsHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In UsersHelp")
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
			Cmd:   "groups addadmins",
			Usage: "groups addadmins <group(s)> <user(s)>",
			Short: "Add user(s) to a group(s) admin list. (groups and users as comma separated list)",
		},
		HelpFmt{
			Cmd:   "groups deladmins",
			Usage: "groups deladmins <group(s)> <user(s)>",
			Short: "Del user(s) from a group(s) admin list. (groups and users as comma separated list)",
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
	gs := lib.ListGroups()
	for _, g := range gs {
		i.Sendf(command.Target, "Group: %s", g.String())
	}

}

func GroupsAdd(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Adding Group(s): %+v", command.Args)
	lib.AddGroups(strings.Split(command.Args[0], ","))
}

func GroupsDel(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Deleting Group(s): %+v", command.Args)
	groups := strings.Split(command.Args[0], ",")
	if fnd, _ := lib.In("Admin", groups); fnd {
		efmt := "User (%s) is not authorized to delete Admin group."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	lib.DelGroups(groups)
}

func GroupsAddMembers(i *Irc, command *Command) {
	if len(command.Args) == 2 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Adding Members(s) to Group(s): %+v to %+v", command.Args[1], command.Args[0])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")
	lib.AddUsersToGroups("Member", gs, us)
}

func GroupsDelMembers(i *Irc, command *Command) {
	if len(command.Args) == 2 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Deleting User(s) from Group(s): %+v to %+v", command.Args[0], command.Args[1])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")
	lib.DelUsersFromGroups("Member", gs, us)
}

func GroupsAddAdmins(i *Irc, command *Command) {
	if len(command.Args) == 2 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Adding Admin(s) to Group(s): %+v to %+v", command.Args[1], command.Args[0])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")
	lib.AddUsersToGroups("Admin", gs, us)
}

func GroupsDelAdmins(i *Irc, command *Command) {
	if len(command.Args) == 2 {
		GroupsHelp(i, command)
		return
	}
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
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
	jww.DEBUG.Printf("Deleting User(s) from Group(s): %+v to %+v", command.Args[0], command.Args[1])
	gs := strings.Split(command.Args[0], ",")
	us := strings.Split(command.Args[1], ",")
	lib.DelUsersFromGroups("Admin", gs, us)
}

func Groups(i *Irc, command *Command) {
	jww.DEBUG.Println("In Users")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("groups", GroupsBase)
	c.HandleFunc("add", GroupsAdd)
	c.HandleFunc("del", GroupsDel)
	c.HandleFunc("addadmins", GroupsAddAdmins)
	c.HandleFunc("deladmins", GroupsDelAdmins)
	c.HandleFunc("addmembers", GroupsAddMembers)
	c.HandleFunc("delmembers", GroupsDelMembers)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
