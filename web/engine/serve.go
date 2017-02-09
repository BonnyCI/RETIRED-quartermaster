package engine

import (
	"io"
	"net"
	"net/http"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

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
