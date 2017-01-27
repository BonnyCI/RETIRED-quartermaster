package database

import (
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	jww "github.com/spf13/jwalterweatherman"
)

// Stop being so dense. Make it non-generic to get it done

type Database struct {
	Db    string
	DbObj *storm.DB
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

type DataS interface {
	Save(*Database)
	Delete(*Database)
	Update(*Database)
	UpdateField(*Database, string, string)
	GetOne(*Database, string, string)
}

func GetAll(db *Database, d interface{}) {
	jww.DEBUG.Printf("Getting all: %T", d)
	if err := db.DbObj.All(d); err != nil {
		jww.ERROR.Printf("Failure to get all: %+v", d)
	}
}

func GetAllByIndex(db *Database, field string, d interface{}) {
	jww.DEBUG.Printf("Getting all: %T (i:%s)", d, field)
	if err := db.DbObj.AllByIndex(field, d); err != nil {
		jww.ERROR.Printf("Failure to get by indexed field: %s", field)
	}
}
