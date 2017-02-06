package bot

import (
	"crypto/tls"
	"strconv"
	"strings"
	"sync"

	irc "github.com/fluffle/goirc/client"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/helpers"
)

type Irc struct {
	Conf     *irc.Config
	Conn     *irc.Conn
	commands *Commands
	help     *Commands
	quit     chan bool
	in       chan string
}

var instance *Irc
var io sync.Once

func GetIrc() *Irc {
	io.Do(func() {
		instance = &Irc{}
	})
	return instance
}

func Configure(i *Irc) {
	server := viper.GetString("server")
	password := viper.GetString("password")
	port := viper.GetString("port")
	usetls, _ := strconv.ParseBool(viper.GetString("UseTLS"))
	nick := viper.GetString("nick")
	user := viper.GetString("user")

	i.Conf = irc.NewConfig(nick, user, user)

	i.Conf.Version = helpers.QuartermasterVersion()
	i.Conf.QuitMessage = "Arrrrrgh!"
	i.Conf.Server = string(server + ":" + port)
	i.Conf.Pass = password
	i.Conf.SSL = usetls
	i.Conf.SSLConfig = &tls.Config{InsecureSkipVerify: true}
	i.Conn = irc.Client(i.Conf)
	i.quit = make(chan bool)
	i.in = make(chan string, 4)
}

func (i *Irc) addCommands() {
	i.commands = &Commands{Handlers: HandlerSet()}
	i.help = &Commands{Handlers: HandlerSet()}

	i.commands.HandleFunc("", Help)
	i.commands.HandleFunc("quit", Quit) //Needs admin perms enforcement

	AddFunc(i.commands, i.help, "help", Help, HelpHelp)
	AddFunc(i.commands, i.help, "notify", Notify, NotifyHelp)
	AddFunc(i.commands, i.help, "ping", Ping, PingHelp)
	AddFunc(i.commands, i.help, "status", Status, StatusHelp)
	AddFunc(i.commands, i.help, "users", Users, UsersHelp)
	AddFunc(i.commands, i.help, "groups", Groups, GroupsHelp)
}

func (i *Irc) Connect() {
	i.addCommands()
	i.Conn.EnableStateTracking()
	i.Conn.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		channels := strings.Split(viper.GetString("channels"), ",")

		for i := range channels {
			jww.DEBUG.Println("Joining Channel: " + channels[i])
			conn.Join(channels[i])
		}
	})

	i.Conn.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		i.Disconnect()
	})

	i.Conn.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		jww.DEBUG.Println("Received message:", line.Target(), line.Text())
		nick := viper.GetString("nick")
		if line.Target()[0] == '#' {
			if !strings.HasPrefix(line.Text(), "!"+nick) {
				return
			}
		}

		c := NewCommand(line.Target(), line.Nick, line.Text())
		i.commands.Handlers.Dispatch(i, c)
	})

	if err := i.Conn.Connect(); err != nil {
		jww.ERROR.Printf("Connection error: %s\n", err)
		return
	}

	<-i.quit
}

func (i *Irc) Disconnect() {
	database.CloseInstance()
	jww.INFO.Println("Disconnecting from IRC Server")
	i.Conn.Quit("Quatermaster is walking the plank!")
	i.quit <- true
}

func (i *Irc) Send(target string, msg string) {
	i.Sendf(target, "%s", msg)
}

func (i *Irc) Sendf(target string, format string, msgs ...string) {
	msgsI := make([]interface{}, len(msgs))
	for k, v := range msgs {
		msgsI[k] = v
	}
	i.Conn.Privmsgf(target, format, msgsI...)
}
