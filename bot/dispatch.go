package bot

import (
	"strings"
	"sync"

	"github.com/fluffle/goirc/logging"
)

type Handler interface {
	Handle(*Irc, *Command)
}

type Remover interface {
	Remove()
}

type cSet struct {
	set map[string]*cList
	sync.RWMutex
}

type cList struct {
	start, end *cNode
}

type cNode struct {
	next, prev *cNode
	set        *cSet
	command    string
	handler    Handler
}

type Commands struct {
	Handlers *cSet
}

type HandlerFunc func(*Irc, *Command)

func (hf HandlerFunc) Handle(i *Irc, command *Command) {
	hf(i, command)
}

func (cn *cNode) Handle(i *Irc, command *Command) {
	cn.handler.Handle(i, command)
}

func (cn *cNode) Remove() {
	cn.set.remove(cn)
}

func HandlerSet() *cSet {
	return &cSet{set: make(map[string]*cList)}
}

func (cs *cSet) add(cmd string, h Handler) Remover {
	cs.Lock()
	defer cs.Unlock()
	cmd = strings.ToLower(cmd)
	l, ok := cs.set[cmd]
	if !ok {
		l = &cList{}
	}
	cn := &cNode{
		set:     cs,
		command: cmd,
		handler: h,
	}
	if !ok {
		l.start = cn
	} else {
		cn.prev = l.end
		l.end.next = cn
	}
	l.end = cn
	cs.set[cmd] = l
	return cn
}

func (cs *cSet) remove(cn *cNode) {
	cs.Lock()
	defer cs.Unlock()
	l, ok := cs.set[cn.command]
	if !ok {
		logging.Error("Removing node for unknown command'%s'", cn.command)
		return
	}
	if cn.next == nil {
		l.end = cn.prev
	} else {
		cn.next.prev = cn.prev
	}
	if cn.prev == nil {
		l.start = cn.next
	} else {
		cn.prev.next = cn.next
	}
	cn.next = nil
	cn.prev = nil
	cn.set = nil
	if l.start == nil || l.end == nil {
		delete(cs.set, cn.command)
	}
}

func (cs *cSet) Dispatch(i *Irc, command *Command) {
	cs.RLock()
	defer cs.RUnlock()
	cmd := strings.ToLower(command.Cmd)
	list, ok := cs.set[cmd]
	if !ok {
		str := cmd
		if command.Parent != nil {
			str = command.Parent.Cmd + " " + cmd
		}

		i.Sendf(command.Sender, "Unknown Command: %s", str)

		return
	}
	wg := &sync.WaitGroup{}
	for cn := list.start; cn != nil; cn = cn.next {
		wg.Add(1)
		go func(cn *cNode) {
			cn.Handle(i, command)
			wg.Done()
		}(cn)
	}
	wg.Wait()
}

func (c *Commands) HandleFunc(name string, hf HandlerFunc) Remover {
	return c.Handle(name, hf)
}

func (c *Commands) Handle(name string, h Handler) Remover {
	return c.Handlers.add(name, h)
}
