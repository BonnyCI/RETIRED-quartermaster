package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/lib"
	"github.com/bonnyci/quartermaster/web/engine"
	"github.com/bonnyci/quartermaster/web/middleware"
)

func StatusUserAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()

	type ApiIn struct {
		Data []string `json:"data"`
	}

	var in ApiIn

	if err := lib.Build(r.Body, &in); err != nil {
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
	date := params["date"]
	index, _ := strconv.Atoi(params["index"])

	user, _, _ := r.BasicAuth()

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

type StatusAPI struct {
	engine.API
}

func GetApi() *StatusAPI {
	return &StatusAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/status/": []engine.HandlersS{
					engine.MakeHandler("GET", StatusBucketListHandleFunc),
					engine.MakeHandler("POST", StatusUserAddHandleFunc, middleware.AuthMiddleware)},
				"/status/{date}":         []engine.HandlersS{engine.MakeHandler("GET", StatusAllGetHandleFunc)},
				"/status/{user}/{date}":  []engine.HandlersS{engine.MakeHandler("GET", StatusUserGetHandleFunc)},
				"/status/{date}/{index}": []engine.HandlersS{engine.MakeHandler("DELETE", StatusUserDelHandleFunc, middleware.AuthMiddleware)},
			},
		},
	}
}
