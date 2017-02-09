package engine

import (
	"net/http"
)

type HandlersS struct {
	Method  string
	Handler http.Handler
}

type HandlersT map[string][]HandlersS

func chain(h http.HandlerFunc, in ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	if len(in) == 0 {
		return h
	}

	c := in[0]

	if len(in) > 1 {
		in = append(in[:0], in[1:]...)
	} else {
		in = nil
	}

	return chain(c(h), in...)

}

func reverse(in ...func(http.HandlerFunc) http.HandlerFunc) []func(http.HandlerFunc) http.HandlerFunc {

	for left, right := 0, len(in)-1; left < right; left, right = left+1, right-1 {
		in[left], in[right] = in[right], in[left]
	}

	return in
}

func MakeHandler(method string, handler http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) HandlersS {
	return HandlersS{
		Method:  method,
		Handler: chain(handler, reverse(middleware...)...),
	}
}
