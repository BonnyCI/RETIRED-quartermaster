package bot

import jww "github.com/spf13/jwalterweatherman"

func UsersHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In UsersHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "users",
			Usage: "users <user|users|group> <action>",
			Short: "Alert a user, list of users (comma seperated), or group of needed action.",
		},
		HelpFmt{
			Cmd:   "users set",
			Usage: "users set <action> <cron> <user|users|group>",
			Short: "Set a cron (std cron format) for a user, list of users (comma seperated), or group for an action.",
		},
		HelpFmt{
			Cmd:   "users del",
			Usage: "users del <action> <user|users|group>",
			Short: "Del a cron for a user, list of users (comma seperated), or group for a given action action.",
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func UsersBase(i *Irc, command *Command) {
}

func Users(i *Irc, command *Command) {
	jww.DEBUG.Println("In Users")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("Users", UsersBase)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
