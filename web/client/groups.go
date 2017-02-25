package client

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/endpoints/auth"
	"github.com/bonnyci/quartermaster/web/endpoints/groups"
)

func GetGroups() ([]database.GroupS, error) {
	var d []database.GroupS
	client := NewClient()

	err := client.Get("http://localhost:8888/groups/", nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}

	return d, nil
}

func GetGroup(g string) (database.GroupS, error) {
	var d database.GroupS
	client := NewClient()

	err := client.Get("http://localhost:8888/groups/"+g, nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	return d, nil
}

func AddGroup(user string, password string, add string) (database.GroupS, error) {
	var d database.GroupS
	client := NewClient()

	in := groups.GroupApiIn{Group: add}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	err = client.Put(token, "http://127.0.0.1:8888/groups/", &in, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	return d, nil
}

func DelGroup(user string, password string, del string) error {
	client := NewClient()

	in := groups.GroupApiIn{Group: del}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	err = client.Delete(token, "http://127.0.0.1:8888/groups/", &in)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	return nil
}

func AddMemberToGroup(user string, password string, group string, add string) (database.GroupS, error) {
	var d database.GroupS
	client := NewClient()

	in := groups.GroupUserApiIn{Username: add}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	err = client.Put(token, "http://127.0.0.1:8888/groups/"+group, &in, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	return d, nil
}

func DelMemberFromGroup(user string, password string, group string, del string) error {
	client := NewClient()

	in := groups.GroupUserApiIn{Username: del}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	err = client.Delete(token, "http://127.0.0.1:8888/groups/"+group, &in)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	return nil
}
