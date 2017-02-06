package web

import (
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/bonnyci/quartermaster/database"
)

type BackupAPI struct {
	Name     string
	Handlers HandlersT
}

func (u *BackupAPI) Get() HandlersT {
	return u.Handlers
}

func BackupHandleFunc(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance()
	err := db.DbObj.Bolt.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="backup.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var backupAPI = &BackupAPI{
	Name: "backup",
	Handlers: HandlersT{
		"/backup/": []HandlersS{MakeHandler("GET", BackupHandleFunc)},
	},
}
