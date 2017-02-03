package web

import (
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/pschwartz/quartermaster/database"
)

func BackupHandleFunc(w http.ResponseWriter, rew *http.Request) {
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

func BackupHTTP() {
	be := GetEngine("backup")
	be.SetAddr(":12800")
	be.Add("/backup/", BackupHandleFunc)
	be.Start()
}
