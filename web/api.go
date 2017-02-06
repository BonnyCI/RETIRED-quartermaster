package web

import (
	"encoding/json"
	"io"
)

func Build(in io.Reader, out interface{}) error {
	if err := json.NewDecoder(in).Decode(out); err != nil {
		return err
	}
	return nil
}

type API interface {
	Get() HandlersT
}

func AddEndpoints(e *Engine, i API) {
	for k, v := range i.Get() {
		e.Add(k, v)
	}
}

func BackupHTTP() {
	be := GetEngine("backup")
	be.SetAddr(":12800")
	AddEndpoints(be, backupAPI)
	be.Start()
}

func ApiHTTP() {
	be := GetEngine("api")
	be.SetAddr(":8888")
	AddEndpoints(be, groupsApi)
	AddEndpoints(be, notifyApi)
	AddEndpoints(be, usersApi)
	AddEndpoints(be, statusApi)
	be.Start()
}
