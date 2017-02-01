package lib

import (
	"time"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/database"
)

func NewStatus(uS database.UserS, s string) database.StatusS {
	jww.DEBUG.Println("In NewStatus")
	st := database.StatusS{
		User:   uS,
		Date:   time.Now().UTC().Format(DFMT),
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
