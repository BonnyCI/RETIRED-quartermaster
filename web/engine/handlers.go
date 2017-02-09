package engine

import (
	"net/http"
)

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
