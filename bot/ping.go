package bot

import (
	"strconv"

	jww "github.com/spf13/jwalterweatherman"
)

func PingHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In PingHelp")
	help := HelpFmt{
		Cmd:   "ping",
		Usage: "ping <cnt>",
		Short: "Ping the Quartermaster to test.",
	}

	help.Use(i, command)
}

func Ping(i *Irc, command *Command) {
	jww.WARN.Println("In Ping")
	if len(command.Args) > 0 {
		cnt, _ := strconv.Atoi(command.Args[0])
		for x := 0; x < cnt; x++ {
			i.Send(command.Target, "PONG")
		}
	} else {
		i.Send(command.Target, "PONG")
	}
}
