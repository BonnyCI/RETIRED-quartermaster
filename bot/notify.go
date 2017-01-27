package bot

import jww "github.com/spf13/jwalterweatherman"

func NotifyHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In NotifyHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "notify",
			Usage: "notify <user|users|group> <action>",
			Short: "Alert a user, list of users (comma seperated), or group of needed action.",
		},
		HelpFmt{
			Cmd:   "notify set",
			Usage: "notify set <action> <cron> <user|users|group>",
			Short: "Set a cron (std cron format) for a user, list of users (comma seperated), or group for an action.",
		},
		HelpFmt{
			Cmd:   "notify del",
			Usage: "notify del <action> <user|users|group>",
			Short: "Del a cron for a user, list of users (comma seperated), or group for a given action action.",
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func NotifyBase(i *Irc, command *Command) {
	jww.DEBUG.Println("notification base")
	i.Send(command.Target, "NOTIFY")
}

func NotifySet(i *Irc, command *Command) {
	jww.DEBUG.Println("Set User notification")
	i.Send(command.Target, "NOTIFY SET")
}

func NotifyDel(i *Irc, command *Command) {
	jww.DEBUG.Println("Delete User notification")
	i.Send(command.Target, "NOTIFY DEL")
}

func Notify(i *Irc, command *Command) {
	jww.DEBUG.Println("Notifing users")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("notify", NotifyBase)
	c.HandleFunc("set", NotifySet)
	c.HandleFunc("del", NotifyDel)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}

func NotifyCron() {

}
