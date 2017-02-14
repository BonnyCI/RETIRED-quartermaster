package engine

import (
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"
)

var engines = make(map[string]*Engine)
var build = make(map[string]*sync.Once)

func GetEngine(e string) *Engine {
	_, ok := build[e]
	if !ok {
		build[e] = &sync.Once{}
	}

	_, en := engines[e]
	if en {
		return engines[e]
	}

	build[e].Do(func() {
		engines[e] = NewEngine(e)
	})

	return engines[e]
}

type EngineI interface {
	Add(string, http.HandlerFunc)
	Remove(string)
	Start()
	SetAddr(string)
}

type Engine struct {
	name     string
	addr     string
	Handlers HandlersT
	Router   *mux.Router
	Closer   io.Closer
}

func NewEngine(name string) *Engine {
	jww.INFO.Printf("Creating new engine: %s", name)
	m := make(HandlersT)
	e := &Engine{name: name, addr: "", Handlers: m, Router: mux.NewRouter()}
	return e
}

func (e *Engine) Add(p string, f []HandlersS) {
	if _, ok := e.Handlers[p]; ok {
		jww.WARN.Printf("Handler %s already in handler map.", p)
		return
	}

	e.Handlers[p] = f
}

func (e *Engine) Remove(p string) {
	if _, ok := e.Handlers[p]; ok {
		jww.DEBUG.Printf("Removing handler %s from handler map.", p)
		delete(e.Handlers, p)
		return
	}
	jww.ERROR.Printf("Handler %s not in handler map.", p)
}

func (e *Engine) SetAddr(a string) {
	jww.INFO.Printf("Setting address: %s for %s", a, e.name)
	e.addr = a
}

func (e *Engine) Start() {
	jww.DEBUG.Println("Building HTTP Engine.")

	e.Router.StrictSlash(true)

	for k, v := range e.Handlers {
		jww.DEBUG.Printf("Adding Handler %s.", k)

		for _, h := range v {
			e.Router.HandleFunc(k, h.Handler.ServeHTTP).Methods(h.Method)
		}
	}
	jww.INFO.Println("Starting HTTP Engine.")
	e.Closer, _ = ListenAndServeWithClose(e.addr, e.Router)
}
