package groups

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

type GroupApiIn struct {
	Group string `json:"group"`
}

type GroupUserApiIn struct {
	Username string `json:"username"`
}

func GroupsListHandleFunc(w http.ResponseWriter, r *http.Request) {
	g, _ := database.ListGroups()
	json.NewEncoder(w).Encode(g)
}

func GroupsAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	var in GroupApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Group == "" {
		jww.ERROR.Println("Group must be set.")
		http.Error(w, fmt.Errorf("Group must be set.").Error(), http.StatusBadRequest)
		return
	}

	if _, err := database.GetGroup(in.Group); err != nil {
		jww.INFO.Printf("Group: %s, does not exist. Creating.", in.Group)
		database.AddGroups([]string{in.Group})
	}

	g, _ := database.GetGroup(in.Group)
	json.NewEncoder(w).Encode(g)
}

func GroupsDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	var in GroupApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Group == "" {
		jww.ERROR.Println("Group must be set.")
		http.Error(w, fmt.Errorf("Group must be set.").Error(), http.StatusBadRequest)
		return
	}

	if _, err := database.GetGroup(in.Group); err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", in.Group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", in.Group).Error(), http.StatusNotFound)
		return
	}

	jww.INFO.Printf("Deleting Group: %s", in.Group)
	database.DelGroups([]string{in.Group})

	json.NewEncoder(w).Encode(map[string]string{"Action": "Group: " + in.Group + " deleted."})
}

func GroupsGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	g, err := database.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(g)
}

func GroupsAddMembersHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	var in GroupUserApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Username == "" {
		jww.ERROR.Println("Username must be set.")
		http.Error(w, fmt.Errorf("Username must be set.").Error(), http.StatusBadRequest)
		return
	}

	g, err := database.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}
	if _, err := database.GetUser(in.Username); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", in.Username)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", in.Username).Error(), http.StatusNotFound)
		return
	}

	database.AddUsersToGroups([]string{group}, []string{in.Username})
	g, _ = database.GetGroup(group)
	json.NewEncoder(w).Encode(g)
}

func GroupsDelMembersHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	var in GroupUserApiIn
	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	if in.Username == "" {
		jww.ERROR.Println("Username must be set.")
		http.Error(w, fmt.Errorf("Username must be set.").Error(), http.StatusBadRequest)
		return
	}

	g, err := database.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}
	if _, err := database.GetUser(in.Username); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", in.Username)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", in.Username).Error(), http.StatusNotFound)
		return
	}

	database.DelUsersFromGroups([]string{group}, []string{in.Username})
	g, _ = database.GetGroup(group)
	json.NewEncoder(w).Encode(g)
}

type GroupsAPI struct {
	engine.API
}

func GetApi() *GroupsAPI {
	return &GroupsAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/groups/": []engine.HandlersS{
					engine.MakeHandler("GET", GroupsListHandleFunc),
					engine.MakeHandler("PUT", GroupsAddHandleFunc, middleware.AuthAndAdmin...),
					engine.MakeHandler("DELETE", GroupsDelHandleFunc, middleware.AuthAndAdmin...),
				},
				"/groups/{group}": []engine.HandlersS{
					engine.MakeHandler("GET", GroupsGetHandleFunc),
					engine.MakeHandler("PUT", GroupsAddMembersHandleFunc, middleware.AuthAndAdmin...),
					engine.MakeHandler("DELETE", GroupsDelMembersHandleFunc, middleware.AuthAndAdmin...),
				},
			},
		},
	}
}
