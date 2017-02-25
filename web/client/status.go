package client

import (
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/endpoints/auth"
	"github.com/bonnyci/quartermaster/web/endpoints/status"
)

func GetStatusList() ([]string, error) {
	var d []string
	client := NewClient()

	err := client.Get("http://localhost:8888/status/", nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}

	return d, nil
}

func GetStatusAll(date string) ([]database.StatusS, error) {
	var d []database.StatusS
	client := NewClient()

	err := client.Get("http://localhost:8888/status/"+date, nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}

	return d, nil
}

func GetStatus(user string, date string) ([]database.StatusS, error) {
	var d []database.StatusS
	client := NewClient()

	err := client.Get("http://localhost:8888/status/"+user+"/"+date, nil, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}

	return d, nil
}

func AddStatus(user string, password string, add string) (database.GroupS, error) {
	var d database.GroupS
	client := NewClient()

	in := status.StatusApiIn{Data: []string{add}}

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	err = client.Put(token, "http://127.0.0.1:8888/status/", &in, &d)
	if err != nil {
		jww.ERROR.Println(err)
		return d, err
	}

	return d, nil
}

func DelStatus(user string, password string, date string, index string) error {
	client := NewClient()

	token, err := client.Auth("http://127.0.0.1:8888/auth/", auth.AuthApiIn{Username: user, Password: password})
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	err = client.Delete(token, "http://127.0.0.1:8888/status/"+date+"/"+index, nil)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}
	return nil
}
