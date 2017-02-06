package database

import (
	jww "github.com/spf13/jwalterweatherman"
)

type UserS struct {
	ID   int    `json:"id"`
	Nick string `json:"nick" storm:"index,unique"`
}

func (s *UserS) String() string {
	return "Nick(" + s.Nick + ")"
}

func (s *UserS) Save() error {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return nil
}

func (s *UserS) Delete() error {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s *UserS) Update() error {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}
