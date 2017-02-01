package database

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type GroupS struct {
	ID      int
	Name    string `storm:"index,unique"`
	Members []string
	Admins  []string
}

func (s *GroupS) String() string {
	m := strings.Join(s.Members, ",")
	a := strings.Join(s.Admins, ",")
	return "Group(" + s.Name + " A:" + a + " M:" + m + ")"
}

func (s *GroupS) Save() {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
	}
}

func (s *GroupS) Delete() {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
	}
}

func (s *GroupS) Update() {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
	}
}
