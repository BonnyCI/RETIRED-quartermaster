package database

import (
	jww "github.com/spf13/jwalterweatherman"
)

type GroupS struct {
	ID      int
	Name    string `storm:"index,unique"`
	Members []string
}

func (s *GroupS) Save(db *Database) {
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
	}
}

func (s *GroupS) Delete(db *Database) {
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
	}
}

func (s *GroupS) Update(db *Database) {
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
	}
}

func (s *GroupS) UpdateField(db *Database, field string, value string) {
	jww.DEBUG.Printf("Updating: %+v (f:%s=v:%s)", s, field, value)
	if err := db.DbObj.UpdateField(s, field, value); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
	}
}

func (s *GroupS) GetOne(db *Database, field string, value string) {
	jww.DEBUG.Printf("Getting: (f:%s=v:%s)", field, value)
	if err := db.DbObj.One(field, value, s); err != nil {
		jww.ERROR.Printf("Failure to get One: %s - %s", field, value)
	}
}
