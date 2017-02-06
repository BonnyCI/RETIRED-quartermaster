package web

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"

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

// Copyright 2016 stackoverflow.com http://stackoverflow.com/questions/39320025/go-how-to-stop-http-listenandserve

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func ListenAndServeWithClose(addr string, handler http.Handler) (sc io.Closer, err error) {
	var listener net.Listener

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			jww.ERROR.Println("HTTP Server Error - ", err)
		}
	}()

	return listener, nil
}

type HandlersS struct {
	Method  string
	Handler http.HandlerFunc
}

type HandlersT map[string][]HandlersS

func MakeHandler(m string, h http.HandlerFunc) HandlersS {
	return HandlersS{
		Method:  m,
		Handler: h,
	}
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

func NewEngine(n string) *Engine {
	jww.INFO.Printf("Creating new engine: %s", n)
	m := make(HandlersT)
	return &Engine{name: n, addr: "", Handlers: m, Router: mux.NewRouter()}
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
	for k, v := range e.Handlers {
		jww.DEBUG.Printf("Adding Handler %s.", k)

		for _, h := range v {
			e.Router.HandleFunc(k, h.Handler).Methods(h.Method)
		}
	}
	jww.INFO.Println("Starting HTTP Engine.")
	e.Closer, _ = ListenAndServeWithClose(e.addr, e.Router)
}
