package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pschwartz/quartermaster/lib"
)

type StatusAPI struct {
	Name     string
	Handlers HandlersT
}

func (s *StatusAPI) Get() HandlersT {
	return s.Handlers
}

func StatusUserAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println("In ADD")
	params := mux.Vars(r)
	user := params["user"]
	jww.DEBUG.Println(user)

	type ApiIn struct {
		Data []string `json:"data"`
	}

	var in ApiIn

	if err := Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	u, err := lib.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	for k := range in.Data {
		lib.NewStatus(u, in.Data[k])
	}

	st := lib.GetStatus(u, lib.DStamp)

	json.NewEncoder(w).Encode(st)
}

func StatusUserDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	date := params["date"]
	index, _ := strconv.Atoi(params["index"])

	jww.DEBUG.Println("In Del, index:", index)

	u, err := lib.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	st := lib.GetStatus(u, date)
	if len(st) < index {
		jww.ERROR.Printf("User: %s has no Status for %s", user, date)
		http.Error(w, fmt.Errorf("User: %s has no Status for %s", user, date).Error(), http.StatusNotFound)
		return
	}

	lib.DelStatus(st[index-1])
	json.NewEncoder(w).Encode(map[string]string{"Action": "User: " + user + ", status index: " + strconv.Itoa(index) + " on " + date + " deleted."})

}

func StatusUserGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	date := params["date"]

	u, err := lib.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	st := lib.GetStatus(u, date)
	if len(st) == 0 {
		jww.ERROR.Printf("User: %s has no Status for %s", user, date)
		http.Error(w, fmt.Errorf("User: %s has no Status for %s", user, date).Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(st)

}

func StatusAllGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println("In ALL")
	params := mux.Vars(r)
	date := params["date"]
	status, err := lib.GetAllStatusByDate(date)
	if err != nil {
		jww.ERROR.Printf("Cannot get all status' for %s", date)
		http.Error(w, fmt.Errorf("Cannot get all status' for %s", date).Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func StatusBucketListHandleFunc(w http.ResponseWriter, r *http.Request) {
	buckets := lib.GetAllStatusBuckets()
	json.NewEncoder(w).Encode(buckets)
}

var statusApi = &StatusAPI{
	Name: "status",
	Handlers: HandlersT{
		"/status/":                      []HandlersS{MakeHandler("GET", StatusBucketListHandleFunc)},
		"/status/{user}":                []HandlersS{MakeHandler("POST", StatusUserAddHandleFunc)},
		"/status/{date}":                []HandlersS{MakeHandler("GET", StatusAllGetHandleFunc)},
		"/status/{user}/{date}":         []HandlersS{MakeHandler("GET", StatusUserGetHandleFunc)},
		"/status/{user}/{date}/{index}": []HandlersS{MakeHandler("DELETE", StatusUserDelHandleFunc)},
	},
}
