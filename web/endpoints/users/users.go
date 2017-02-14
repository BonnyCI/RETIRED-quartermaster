package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/engine"
	"github.com/bonnyci/quartermaster/web/middleware"
)

type UserApiIn struct {
	Username string `json:"user"`
	Password string `json:"password,omitempty"`
}

func UsersListHandleFunc(w http.ResponseWriter, r *http.Request) {
	u := database.ListUsers()
	json.NewEncoder(w).Encode(u)
}

func UsersAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	var in UserApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Username == "" {
		jww.ERROR.Println("Username must be set.")
		http.Error(w, fmt.Errorf("Username must be set.").Error(), http.StatusBadRequest)
		return
	}

	if _, err := database.GetUser(in.Username); err != nil {
		jww.INFO.Printf("User: %s, does not exist. Creating.", in.Username)
		u := database.UserS{Nick: in.Username, Password: in.Password}
		u.Save()
	}

	u, _ := database.GetUser(in.Username)
	json.NewEncoder(w).Encode(u)
}

func UsersGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	u, err := database.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func UsersDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	var in UserApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if _, err := database.GetUser(in.Username); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", in.Username)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", in.Username).Error(), http.StatusNotFound)
		return
	}

	jww.INFO.Printf("Deleting User: %s", in.Username)
	database.DelUsers([]string{in.Username})

	json.NewEncoder(w).Encode(map[string]string{"Action": "User: " + in.Username + " deleted."})
}

type UsersAPI struct {
	engine.API
}

func GetApi() *UsersAPI {
	return &UsersAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/users/": []engine.HandlersS{
					engine.MakeHandler("GET", UsersListHandleFunc),
					engine.MakeHandler("PUT", UsersAddHandleFunc, middleware.AuthAndAdmin...),
					engine.MakeHandler("DELETE", UsersDelHandleFunc, middleware.AuthAndAdmin...),
				},
				"/users/{user}": []engine.HandlersS{engine.MakeHandler("GET", UsersGetHandleFunc)},
			},
		},
	}
}
