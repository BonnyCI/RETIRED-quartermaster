package backup

import (
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/engine"
	"github.com/bonnyci/quartermaster/web/middleware"
)

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

type BackupAPI struct {
	engine.API
}

func GetApi() *BackupAPI {
	return &BackupAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/backup/": []engine.HandlersS{engine.MakeHandler("GET", BackupHandleFunc, middleware.AdminMiddleware, middleware.AuthMiddleware)},
			},
		},
	}
}
