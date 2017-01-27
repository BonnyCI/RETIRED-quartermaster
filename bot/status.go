package bot

import jww "github.com/spf13/jwalterweatherman"

func StatusHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In StatusHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "status",
			Usage: "status <text>",
			Short: "Add a status line to the standup record for the irc user",
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func StatusBase(i *Irc, command *Command) {
	jww.DEBUG.Println("In Status")
	i.Send(command.Sender, "STATUS")
}

func Status(i *Irc, command *Command) {
	jww.DEBUG.Println("In Status")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("status", StatusBase)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
