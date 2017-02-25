package database

import (
	jww "github.com/spf13/jwalterweatherman"
	"time"
)

const DFMT = "2006-01-02"

var DStamp = time.Now().UTC().Format(DFMT)

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

func NewStatus(uS UserS, s string) StatusS {
	jww.DEBUG.Println("In NewStatus")
	st := StatusS{
		User:   uS,
		Date:   DStamp,
		Status: s,
	}
	st.Save()
	return st
}

func GetStatus(uS UserS, date string) []StatusS {
	jww.DEBUG.Printf("Getting status for %s on %s", uS.Nick, date)
	var st []StatusS
	status := DateBucket("status", date)
	status.Find("User", uS, &st)
	return st
}

func DelStatus(sT StatusS) {
	jww.DEBUG.Printf("Deleting status for %s from %s.", sT.User, sT.Date)
	sT.Delete()
}

func GetAllStatusBuckets() []string {
	jww.DEBUG.Println("Getting list of all status buckets.")
	return BucketList("status")
}

func GetAllStatusByDate(date string) ([]StatusS, error) {
	var d []StatusS
	status := DateBucket("status", date)

	jww.DEBUG.Printf("Getting all: %T", d)
	if err := status.All(&d); err != nil {
		jww.ERROR.Printf("Failure to get all: %+v", d)
		return d, err
	}

	return d, nil
}
