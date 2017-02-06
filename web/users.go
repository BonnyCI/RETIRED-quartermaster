package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

type UsersAPI struct {
	Name     string
	Handlers HandlersT
}

func (u *UsersAPI) Get() HandlersT {
	return u.Handlers
}

func UsersListHandleFunc(w http.ResponseWriter, r *http.Request) {
	u := lib.ListUsers()
	json.NewEncoder(w).Encode(u)
}

func UsersAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	if _, err := lib.GetUser(user); err != nil {
		jww.INFO.Printf("User: %s, does not exist. Creating.", user)
		lib.AddUsers([]string{user})
	}

	u, _ := lib.GetUser(user)
	json.NewEncoder(w).Encode(u)
}

func UsersGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	u, err := lib.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func UsersDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	if _, err := lib.GetUser(user); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	jww.INFO.Printf("Deleting User: %s", user)
	lib.DelUsers([]string{user})

	json.NewEncoder(w).Encode(map[string]string{"Action": "User: " + user + " deleted."})
}

var usersApi = &UsersAPI{
	Name: "users",
	Handlers: HandlersT{
		"/users/": []HandlersS{MakeHandler("GET", UsersListHandleFunc)},
		"/users/{user}": []HandlersS{
			MakeHandler("PUT", UsersAddHandleFunc),
			MakeHandler("DELETE", UsersDelHandleFunc),
			MakeHandler("GET", UsersGetHandleFunc),
		},
	},
}
