package lib

import (
	"net/http"
	"strconv"

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
	u, _ := GetUser(us)
	u.Delete()
}

func DelUsers(us []string) {
	for _, u := range us {
		delUser(u)
	}
}

func AuthToken(us string, p string, r *http.Request) bool {
	u, err := GetUser(us)

	if err != nil {
		jww.ERROR.Printf("Cannot load user: %s", us)
		return false
	}

	return p == u.Token
}

func IsAdmin(user string) bool {
	u, err := GetUser(user)
	if err != nil {
		return false
	}

	return UserInGroup("Admin", u)
}

func IsSelf(me, user string) bool {
	um, err := GetUser(me)
	if err != nil {
		return false
	}

	u, err := GetUser(user)
	if err != nil {
		return false
	}

	jww.DEBUG.Printf("%s == %s : %s", u, um, strconv.FormatBool(u == um))
	return u == um
}
