package lib

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
)

func NewStatus(uS database.UserS, s string) database.StatusS {
	jww.DEBUG.Println("In NewStatus")
	st := database.StatusS{
		User:   uS,
		Date:   DStamp,
		Status: s,
	}
	st.Save()
	return st
}

func GetStatus(uS database.UserS, date string) []database.StatusS {
	jww.DEBUG.Printf("Getting status for %s on %s", uS.Nick, date)
	var st []database.StatusS
	status := database.DateBucket("status", date)
	status.Find("User", uS, &st)
	return st
}

func DelStatus(sT database.StatusS) {
	jww.DEBUG.Printf("Deleting status for %s from %s.", sT.User, sT.Date)
	sT.Delete()
}

func GetAllStatusBuckets() []string {
	jww.DEBUG.Println("Getting list of all status buckets.")
	return database.BucketList("status")
}

func GetAllStatusByDate(date string) ([]database.StatusS, error) {
	var d []database.StatusS
	status := database.DateBucket("status", date)

	jww.DEBUG.Printf("Getting all: %T", d)
	if err := status.All(&d); err != nil {
		jww.ERROR.Printf("Failure to get all: %+v", d)
		return d, err
	}

	return d, nil
}
