package database

import (
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type GroupS struct {
	ID      int     `json:"id"`
	Name    string  `json:"name" storm:"index,unique"`
	Members []UserS `json:"members,omitempty"`
}

func (s GroupS) String() string {
	var ms []string

	for _, v := range s.Members {
		ms = append(ms, v.Nick)
	}

	m := strings.Join(ms, ",")

	return "(" + s.Name + " M:" + m + ")"
}

func (s GroupS) Compare(d DataS) bool {
	c := d.(GroupS)

	if s.Name != c.Name {
		return false
	}

	if s.Members == nil && c.Members == nil {
		return true
	}

	if s.Members == nil || c.Members == nil {
		return false
	}

	if len(s.Members) != len(c.Members) {
		return false
	}

	for i := range s.Members {
		if !s.Members[i].Compare(c.Members[i]) {
			return false
		}
	}

	return true
}

func (s GroupS) Save() error {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return nil
}

func (s GroupS) Delete() error {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s GroupS) Update() error {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}
