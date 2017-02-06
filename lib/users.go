package lib

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
)

func ListUsers() []database.UserS {
	var du []database.UserS
	database.GetAll(&du)
	return du
}

func GetUser(n string) (database.UserS, error) {
	var u database.UserS
	if err := database.One("Nick", n, &u); err != nil {
		return u, err
	}
	return u, nil
}

func addUser(us string) {
	jww.DEBUG.Printf("Adding user: %s", us)
	u := database.UserS{Nick: us}
	u.Save()
}

func AddUsers(us []string) {
	for _, u := range us {
		addUser(u)
	}
}

func delUser(us string) {
	jww.DEBUG.Printf("Deleting user: %s", us)
	//var u database.UserS
	u, _ := GetUser(us)
	u.Delete()
}

func DelUsers(us []string) {
	for _, u := range us {
		delUser(u)
	}
}
