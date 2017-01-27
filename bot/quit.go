package bot

import jww "github.com/spf13/jwalterweatherman"

func Quit(i *Irc, command *Command) {
	jww.DEBUG.Println("Quiting IRC")
	i.Disconnect()
}
