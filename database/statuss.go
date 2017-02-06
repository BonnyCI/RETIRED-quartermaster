package database

import (
	jww "github.com/spf13/jwalterweatherman"
)

type StatusS struct {
	ID     int    `json:"id"`
	Date   string `json:"date" storm:"index"`
	Status string `json:"status"`
	User   UserS  `json:"user" storm:"index"`
}

func (s *StatusS) String() string {
	return "Status(U:" + s.User.Nick + " C:" + s.Date + ")"
}

func (s *StatusS) Save() error {
	jww.DEBUG.Printf("Saving: %+v", s)
	bucket := DateBucket("status", s.Date)
	if err := bucket.Save(s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return nil
}

func (s *StatusS) Delete() error {
	bucket := DateBucket("status", s.Date)
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := bucket.DeleteStruct(s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s *StatusS) Update() error {
	bucket := DateBucket("status", s.Date)
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := bucket.Update(s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}
