package engine

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
