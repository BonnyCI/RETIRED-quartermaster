package engine

import (
	"encoding/json"
	"io"
)

type API interface {
	Get() HandlersT
	AddToEngine(*Engine)
}

type APIBase struct {
	Handlers HandlersT
}

func (a APIBase) Get() HandlersT {
	return a.Handlers
}

func (a APIBase) AddToEngine(e *Engine) {
	for k, v := range a.Get() {
		e.Add(k, v)
	}
}

func Build(in io.Reader, out interface{}) error {
	if err := json.NewDecoder(in).Decode(out); err != nil {
		return err
	}
	return nil
}
