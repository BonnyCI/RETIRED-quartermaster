package database

import (
	"github.com/pborman/uuid"
	jww "github.com/spf13/jwalterweatherman"
)

type UserS struct {
	ID    int    `json:"id"`
	Nick  string `json:"nick" storm:"index,unique"`
	Token string `json:"-" storm:"unique"`
}

func (s UserS) String() string {
	return "Nick(" + s.Nick + ")"
}

func (s UserS) Compare(d DataS) bool {
	c := d.(UserS)

	if s == c {
		return true
	}
	return false
}

func (s UserS) Save() error {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return s.GenerateToken()
}

func (s UserS) Delete() error {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s UserS) Update() error {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}

func (s UserS) GenerateToken() error {
	jww.DEBUG.Printf("Generating Token for %s", s.Nick)
	jww.DEBUG.Printf("Token1 %s", s.Token)
	if s.Token == "" {
		s.Token = uuid.New()
		jww.DEBUG.Printf("Token2 %s", s.Token)
		return s.Update()
	}
	return nil
}
