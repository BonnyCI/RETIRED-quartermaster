package bot

import (
	"strings"
)

type Command struct {
	Parent *Command
	Cmd    string
	Target string
	Sender string
	Nick   string
	Args   []string
}

func (c *Command) GetSubCommand() *Command {
	var cmd []string
	if c.Target[0] == '#' {
		cmd = append([]string{"!" + c.Nick}, c.Args...)
	} else {
		cmd = c.Args
	}

	cm := NewCommand(c.Target, c.Sender, strings.Join(cmd, " "), c.Nick)
	cm.Parent = c
	return cm
}

func (c *Command) SplitLine(target string, sender string, cli string) {
	line := strings.Split(cli, " ")
	c.Target = target
	c.Sender = sender
	switch target[0] {
	case '#':
		switch line[0][0] {
		case '&', '+', '!':
			if strings.HasPrefix(line[0], "!"+c.Nick) {
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

func NewCommand(target string, sender string, cli string, nick string) *Command {
	c := &Command{Nick: nick}
	c.SplitLine(target, sender, cli)
	return c
}
