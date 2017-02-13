package database

import (
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	jww "github.com/spf13/jwalterweatherman"
)

type DataS interface {
	Save() error
	Delete() error
	Update() error
	Compare(DataS) bool
}

func InitStatus() {
	db := GetInstance()
	db.DbObj.Bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("status"))
		if err != nil {
			jww.DEBUG.Println("Could not create status bucket")
		}
		return err
	})
}

func GetAll(d interface{}) error {
	db := GetInstance()
	jww.DEBUG.Printf("Getting all: %T", d)
	if err := db.DbObj.All(d); err != nil {
		jww.ERROR.Printf("Failure to get all: %+v", d)
		return err
	}
	return nil
}

func GetAllByIndex(field string, d interface{}) error {
	db := GetInstance()
	jww.DEBUG.Printf("Getting all: %T (i:%s)", d, field)
	if err := db.DbObj.AllByIndex(field, d); err != nil {
		jww.ERROR.Printf("Failure to get by indexed field: %s", field)
		return err
	}
	return nil
}

func One(field string, value interface{}, to interface{}) error {
	db := GetInstance()
	jww.DEBUG.Printf("Getting: (f:%s=v:%s)", field, value)
	if err := db.DbObj.One(field, value, to); err != nil {
		jww.ERROR.Printf("Failure to find: %s - %s", field, value)
		return err
	}
	return nil
}

func Find(field string, value interface{}, to interface{}) error {
	db := GetInstance()
	jww.DEBUG.Printf("Getting: (f:%s=v:%s)", field, value)
	if err := db.DbObj.Find(field, value, to); err != nil {
		jww.ERROR.Printf("Failure to find: %s - %s", field, value)
		return err
	}
	return nil
}

func DateBucket(p string, d string) storm.Node {
	db := GetInstance()
	status := db.DbObj.From(p)
	node := status.From(d)
	return node
}

// There is no way with storm to get a list of all buckets. This uses a direct Bolt view to get the list
func BucketList(p string) []string {
	db := GetInstance()
	InitStatus()

	buckets := []string{}
	db.DbObj.Bolt.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(p)).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			buckets = append(buckets, string(k))
		}
		return nil
	})

	return buckets
}
