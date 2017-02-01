package bot

import (
	"strconv"
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

func StatusHelp(i *Irc, command *Command) {
	jww.DEBUG.Println("In StatusHelp")
	help := []HelpFmt{
		HelpFmt{
			Cmd:   "status add",
			Usage: "status add <text>",
			Short: "Add a status line to the standup record for the irc user",
		},
		HelpFmt{
			Cmd:   "status get",
			Usage: "status get <date>",
			Short: "Get your users status from the given date. (Format: YYYY-MM-DD, default: " + lib.DStamp,
		},
		HelpFmt{
			Cmd:   "status del",
			Usage: "status del <index> <date>",
			Short: "Delete your users status at index from the given date. (Format: YYYY-MM-DD, default: " + lib.DStamp,
		},
	}

	for k := range help {
		help[k].Use(i, command)
	}
}

func StatusAdd(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		StatusHelp(i, command)
		return
	}

	jww.DEBUG.Printf("Adding status for %s.", command.Sender)
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}

	sT := lib.NewStatus(u, strings.Join(command.Args, " "))
	i.Sendf(command.Sender, "Status for %s recieved.", sT.Date)
}

func StatusGet(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		command.Args = append(command.Args, lib.DStamp)
	}
	jww.DEBUG.Printf("Getting status for %s on %s", command.Sender, command.Args[0])
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}
	st := lib.GetStatus(u, command.Args[0])
	i.Sendf(command.Sender, "Status for %s:", command.Args[0])
	if len(st) == 0 {
		i.Send(command.Sender, "No status for this date.")
		return
	}

	for k, v := range st {
		i.Sendf(command.Sender, "#%s: %s", strconv.Itoa(k+1), v.Status)
	}
}

func StatusDel(i *Irc, command *Command) {
	if len(command.Args) == 0 {
		StatusHelp(i, command)
		return
	}
	if len(command.Args) == 1 {
		command.Args = append(command.Args, lib.DStamp)
	}
	jww.DEBUG.Printf("Deleting status for %s at index %s for %s", command.Sender, command.Args[0], command.Args[1])
	index, _ := strconv.Atoi(command.Args[0])
	u, err := lib.GetUser(command.Sender)
	if err != nil {
		efmt := "User (%s) is not registered."
		jww.ERROR.Printf(efmt, command.Sender)
		i.Sendf(command.Target, efmt, command.Sender)
		return
	}
	st := lib.GetStatus(u, command.Args[1])

	if len(st) < index {
		i.Send(command.Sender, "No status for this date.")
		return
	}

	lib.DelStatus(st[index-1])
	i.Sendf(command.Sender, "Status %s for %s deleted.", command.Args[0], command.Args[1])
}

func Status(i *Irc, command *Command) {
	jww.DEBUG.Println("In Status")

	c := Commands{Handlers: HandlerSet()}
	c.HandleFunc("add", StatusAdd)
	c.HandleFunc("del", StatusDel)
	c.HandleFunc("get", StatusGet)

	if len(command.Args) == 0 {
		c.Handlers.Dispatch(i, command)
		return
	}
	c.Handlers.Dispatch(i, command.GetSubCommand())
}
