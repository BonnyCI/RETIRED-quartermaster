package database

import (
	"net/http"
	"strconv"

	jww "github.com/spf13/jwalterweatherman"
)

type UserS struct {
	ID       int    `json:"id"`
	Nick     string `json:"nick" storm:"index,unique"`
	Password string `json:"-"`
}

func (s UserS) String() string {
	return "Nick(" + s.Nick + ")"
}

func (s UserS) Compare(d DataS) bool {
	c := d.(UserS)

	if s == c {
		return true
	}
	return false
}

func (s UserS) Save() error {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(&s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return nil
}

func (s UserS) Delete() error {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(&s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s UserS) Update() error {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(&s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}

func ListUsers() []UserS {
	var du []UserS
	GetAll(&du)
	return du
}

func GetUser(n string) (UserS, error) {
	var u UserS
	if err := One("Nick", n, &u); err != nil {
		return u, err
	}
	return u, nil
}

func addUser(us string) {
	jww.DEBUG.Printf("Adding user: %s", us)
	u := UserS{Nick: us}
	u.Save()
}

func AddUsers(us []string) {
	for _, u := range us {
		addUser(u)
	}
}

func AddPassword(us string, pw string) error {
	u, err := GetUser(us)

	if err != nil {
		jww.ERROR.Printf("Cannot load user: %s", us)
		return err
	}

	u.Password = pw
	return u.Update()
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

func AuthUser(us string, p string, r *http.Request) bool {
	u, err := GetUser(us)

	if err != nil {
		jww.ERROR.Printf("Cannot load user: %s", us)
		return false
	}

	return p == u.Password
}

func IsAdmin(user string) bool {
	u, err := GetUser(user)
	if err != nil {
		return false
	}

	g, err := GetGroup("Admin")
	if err != nil {
		return false
	}

	return UserInGroup(g, u)
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

func RemoveUser(list []UserS, rm UserS) []UserS {
	for k, v := range list {
		if ok := v.Compare(rm); ok {
			return append(list[:k], list[k+1:]...)
		}
	}
	return list
}
