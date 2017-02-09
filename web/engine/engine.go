package engine

import (
	"io"
	"net/http"
	"sync"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/urfave/negroni"

	"github.com/bonnyci/quartermaster/lib"
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
	GetMiddleware(string, ...*negroni.Negroni) *negroni.Negroni
	GetRouter(string, ...*mux.Router) *mux.Router
	Remove(string)
	Start()
	SetAddr(string)
}

type Engine struct {
	name       string
	addr       string
	Handlers   HandlersT
	Router     map[string]*mux.Router
	Middleware map[string]*negroni.Negroni
	Closer     io.Closer
}

func NewEngine(name string) *Engine {
	jww.INFO.Printf("Creating new engine: %s", name)
	m := make(HandlersT)
	r := make(map[string]*mux.Router)
	n := make(map[string]*negroni.Negroni)
	e := &Engine{name: name, addr: "", Handlers: m, Router: r, Middleware: n}
	e.GetRouter("base")
	e.GetMiddleware("GET")
	e.GetMiddleware("PUT")
	e.GetMiddleware("POST")
	e.GetMiddleware("DELETE")
	return e
}

func (e *Engine) GetRouter(name string, r ...*mux.Router) *mux.Router {
	if er, ok := e.Router[name]; ok {
		return er
	}

	if len(r) == 0 || (len(r) != 0 && r[0] != nil) {
		r = append(r, mux.NewRouter())
	}
	e.Router[name] = r[0]
	return e.Router[name]
}

func (e *Engine) GetMiddleware(name string, m ...*negroni.Negroni) *negroni.Negroni {
	if n, ok := e.Middleware[name]; ok {
		return n
	}

	if len(m) == 0 || (len(m) != 0 && m[0] != nil) {
		m = append(m, negroni.New())
	}

	e.Middleware[name] = m[0]
	return e.Middleware[name]
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

	authOpts := httpauth.AuthOptions{
		Realm:    "quartermaster",
		AuthFunc: lib.AuthToken,
	}

	r := e.GetRouter("base")
	n := e.GetMiddleware("base")

	for k, v := range e.Handlers {
		jww.DEBUG.Printf("Adding Handler %s.", k)

		for _, h := range v {
			if h.Method == "GET" {
				r.HandleFunc(k, h.Handler).Methods(h.Method)
			} else {
				r.HandleFunc(k, h.Handler).Methods(h.Method)
			}
		}
	}
	jww.INFO.Println("Starting HTTP Engine.")
	n.UseHandler(r)
	e.Closer, _ = ListenAndServeWithClose(e.addr, httpauth.BasicAuth(authOpts)(n))
}
