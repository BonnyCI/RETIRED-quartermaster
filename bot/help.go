package bot

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type HelpFmt struct {
	Cmd   string
	Usage string
	Short string
}

func Spaces(cmd string) string {
	return strings.Repeat(" ", len(cmd)+2)
}

func (h *HelpFmt) Use(i *Irc, command *Command) {
	i.Sendf(command.Sender, "Description: %s: %s", h.Cmd, h.Short)
	i.Sendf(command.Sender, "Usage: %s", h.Usage)
}

func HelpHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In HelpHelp")
	help := "!" + i.Conf.Me.Nick + " <cmd> {subcmd} <*args> (Commands:"
	for k := range i.commands.Handlers.set {
		if k == "quit" {
			continue
		}
		help = help + " " + k
	}
	help = help + ")"
	i.Send(command.Sender, help)
}

func Help(i *Irc, command *Command) {
	jww.WARN.Println("In Help")

	if len(command.Args) == 0 {
		i.help.Handlers.Dispatch(i, command)
		return
	}

	var cmd []string

	cmd = command.Args
	if command.Target[0] == '#' {
		cmd = append([]string{}, command.Args[:1]...)
	}

	cm := NewCommand(command.Target, command.Sender, strings.Join(cmd, " "), i.Conf.Me.Nick)
	i.help.Handlers.Dispatch(i, cm)

}
