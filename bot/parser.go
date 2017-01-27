package bot

import (
	"strings"

	"github.com/spf13/viper"
)

type Command struct {
	Parent *Command
	Cmd    string
	Target string
	Sender string
	Args   []string
}

func (c *Command) GetSubCommand() *Command {
	var cmd []string
	if c.Target[0] == '#' {
		nick := viper.GetString("nick")
		cmd = append([]string{"!" + nick}, c.Args...)
	} else {
		cmd = c.Args
	}

	cm := NewCommand(c.Target, c.Sender, strings.Join(cmd, " "))
	cm.Parent = c
	return cm
}

func (c *Command) SplitLine(target string, sender string, cli string) {
	line := strings.Split(cli, " ")
	c.Target = target
	c.Sender = sender
	nick := viper.GetString("nick")
	switch target[0] {
	case '#':
		switch line[0][0] {
		case '&', '+', '!':
			if strings.HasPrefix(line[0], "!"+nick) {
				c.Cmd = line[1]
				c.Args = line[2:]
			} else {
				c.Cmd = ""
				c.Args = nil
			}
		}
	default:
		c.Cmd = line[0]
		c.Args = line[1:]
	}
}

func NewCommand(target string, sender string, cli string) *Command {
	c := &Command{}
	c.SplitLine(target, sender, cli)
	return c
}
