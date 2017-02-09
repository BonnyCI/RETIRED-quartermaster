package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/lib"
	"github.com/bonnyci/quartermaster/web/engine"
	"github.com/bonnyci/quartermaster/web/middleware"
)

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

type UsersAPI struct {
	engine.API
}

func GetApi() *UsersAPI {
	return &UsersAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/users/": []engine.HandlersS{engine.MakeHandler("GET", UsersListHandleFunc)},
				"/users/{user}": []engine.HandlersS{
					engine.MakeHandler("PUT", middleware.AdminMiddleware(UsersAddHandleFunc)),
					engine.MakeHandler("DELETE", middleware.AdminMiddleware(UsersDelHandleFunc)),
					engine.MakeHandler("GET", UsersGetHandleFunc),
				},
			},
		},
	}
}
