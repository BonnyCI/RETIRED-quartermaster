package client

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/endpoints/auth"
	"github.com/bonnyci/quartermaster/web/endpoints/users"
)

func GetUsers() ([]database.UserS, error) {
	var d []database.UserS
	client := NewClient()

	err := client.Get("http://localhost:8888/users/", nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}

	return d, nil
}

func GetUser(u string) (database.UserS, error) {
	var d database.UserS
	client := NewClient()

	err := client.Get("http://127.0.0.1:8888/users/"+u, nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	return d, nil
}

func AddUser(user string, password string, add string) (database.UserS, error) {
	var d database.UserS
	client := NewClient()

	in := users.UserApiIn{Username: add}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	err = client.Put(token, "http://127.0.0.1:8888/users/", &in, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}
	return d, nil
}

func DelUser(user string, password string, del string) error {
	client := NewClient()

	in := users.UserApiIn{Username: del}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	err = client.Delete(token, "http://127.0.0.1:8888/users/", &in)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}
	return nil
}
