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
	jww.DEBUG.Printf("Adding Group(s): %+v", command.Args)
	lib.AddGroups(strings.Split(command.Args[0], ","))
}

func GroupsDel(i *Irc, command *Command) {
	jww.DEBUG.Printf("Deleting Group(s): %+v", command.Args)
	lib.DelGroups(strings.Split(command.Args[0], ","))
}

func GroupsAddMembers(i *Irc, command *Command) {
	jww.DEBUG.Printf("Adding Members(s) to Group(s): %+v to %+v", command.Args[1], command.Args[0])
	g := strings.Split(command.Args[0], ",")
	u := strings.Split(command.Args[1], ",")
	lib.AddUsersToGroups("Member", g, u)
}

func GroupsDelMembers(i *Irc, command *Command) {
	jww.DEBUG.Printf("Deleting User(s) from Group(s): %+v to %+v", command.Args[0], command.Args[1])
	g := strings.Split(command.Args[0], ",")
	u := strings.Split(command.Args[1], ",")
	lib.DelUsersFromGroups("Member", g, u)
}

func GroupsAddAdmins(i *Irc, command *Command) {
	jww.DEBUG.Printf("Adding Admin(s) to Group(s): %+v to %+v", command.Args[1], command.Args[0])
	g := strings.Split(command.Args[0], ",")
	u := strings.Split(command.Args[1], ",")
	lib.AddUsersToGroups("Admin", g, u)
}

func GroupsDelAdmins(i *Irc, command *Command) {
	jww.DEBUG.Printf("Deleting User(s) from Group(s): %+v to %+v", command.Args[0], command.Args[1])
	g := strings.Split(command.Args[0], ",")
	u := strings.Split(command.Args[1], ",")
	lib.DelUsersFromGroups("Admin", g, u)
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
