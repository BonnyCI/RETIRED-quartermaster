package database

import (
	"os"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type Database struct {
	Db    string
	DbObj *storm.DB
}

var instance *Database
var open, close sync.Once

func GetInstance() *Database {
	open.Do(func() {
		db := viper.GetString("Database")
		instance = &Database{Db: db}
		instance.Open()
	})
	return instance
}

func CloseInstance() {
	close.Do(func() {
		instance.Close()
	})
}

func (d *Database) Open() {
	jww.DEBUG.Println("Creating/Opening DB:", d.Db)
	DB, err := storm.Open(d.Db, storm.AutoIncrement(),
		storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))

	if err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}

	d.DbObj = DB
}

func (d *Database) Close() {
	jww.DEBUG.Println("Closing DB:", d.Db)
	if err := d.DbObj.Close(); err != nil {
		jww.ERROR.Println(err)
	}
}
