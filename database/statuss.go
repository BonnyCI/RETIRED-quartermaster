package database

import (
	jww "github.com/spf13/jwalterweatherman"
)

type StatusS struct {
	ID     int
	Date   string `storm:"index"`
	Status string
	User   UserS `storm:"index"`
}

func (s *StatusS) String() string {
	return "Status(U:" + s.User.Nick + " C:" + s.Date + ")"
}

func (s *StatusS) Save() {
	jww.DEBUG.Printf("Saving: %+v", s)
	bucket := DateBucket("status", s.Date)
	if err := bucket.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
	}
}

func (s *StatusS) Delete() {
	bucket := DateBucket("status", s.Date)
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := bucket.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
	}
}

func (s *StatusS) Update() {
	bucket := DateBucket("status", s.Date)
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := bucket.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
	}
}
