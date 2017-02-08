package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/lib"
)

type GroupsAPI struct {
	Name     string
	Handlers HandlersT
}

func (u *GroupsAPI) Get() HandlersT {
	return u.Handlers
}

func GroupsListHandleFunc(w http.ResponseWriter, r *http.Request) {
	g, _ := lib.ListGroups()
	json.NewEncoder(w).Encode(g)
}

func GroupsAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	if _, err := lib.GetGroup(group); err != nil {
		jww.INFO.Printf("Group: %s, does not exist. Creating.", group)
		lib.AddGroups([]string{group})
	}

	g, _ := lib.GetGroup(group)
	json.NewEncoder(w).Encode(g)
}

func GroupsGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	g, err := lib.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(g)
}

func GroupsDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	if _, err := lib.GetGroup(group); err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}

	jww.INFO.Printf("Deleting Group: %s", group)
	lib.DelGroups([]string{group})

	json.NewEncoder(w).Encode(map[string]string{"Action": "Group: " + group + " deleted."})
}

func GroupsAddMembersHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]
	user := params["user"]

	g, err := lib.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}
	if _, err := lib.GetUser(user); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	lib.AddUsersToGroups([]string{group}, []string{user})
	g, _ = lib.GetGroup(group)
	json.NewEncoder(w).Encode(g)
}

func GroupsDelMembersHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]
	user := params["user"]

	g, err := lib.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}
	if _, err := lib.GetUser(user); err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	lib.DelUsersFromGroups([]string{group}, []string{user})
	g, _ = lib.GetGroup(group)
	json.NewEncoder(w).Encode(g)
}

var groupsApi = &GroupsAPI{
	Name: "groups",
	Handlers: HandlersT{
		"/groups/": []HandlersS{MakeHandler("GET", GroupsListHandleFunc)},
		"/groups/{group}": []HandlersS{
			MakeHandler("GET", GroupsGetHandleFunc),
			MakeHandler("POST", AdminMiddleware(GroupsAddHandleFunc)),
			MakeHandler("DELETE", AdminMiddleware(GroupsDelHandleFunc)),
		},
		"/groups/{group}/{user}": []HandlersS{
			MakeHandler("POST", AdminMiddleware(GroupsAddMembersHandleFunc)),
			MakeHandler("DELETE", AdminMiddleware(GroupsDelMembersHandleFunc)),
		},
	},
}
