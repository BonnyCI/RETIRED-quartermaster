package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/engine"
	"github.com/bonnyci/quartermaster/web/middleware"
)

type StatusApiIn struct {
	Data []string `json:"data"`
}

func StatusUserAddHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	var in StatusApiIn

	if err := engine.Build(r.Body, &in); err != nil {
		jww.ERROR.Println(err)
		return
	}

	u, err := database.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	for k := range in.Data {
		database.NewStatus(u, in.Data[k])
	}

	st := database.GetStatus(u, database.DStamp)

	json.NewEncoder(w).Encode(st)
}

func StatusUserDelHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	date := params["date"]
	index, _ := strconv.Atoi(params["index"])

	u, err := database.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	st := database.GetStatus(u, date)
	if len(st) < index {
		jww.ERROR.Printf("User: %s has no Status for %s", user, date)
		http.Error(w, fmt.Errorf("User: %s has no Status for %s", user, date).Error(), http.StatusNotFound)
		return
	}

	database.DelStatus(st[index-1])
	json.NewEncoder(w).Encode(map[string]string{"Action": "User: " + user + ", status index: " + strconv.Itoa(index) + " on " + date + " deleted."})

}

func StatusUserGetHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	date := params["date"]

	u, err := database.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	st := database.GetStatus(u, date)
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
	status, err := database.GetAllStatusByDate(date)
	if err != nil {
		jww.ERROR.Printf("Cannot get all status' for %s", date)
		http.Error(w, fmt.Errorf("Cannot get all status' for %s", date).Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func StatusBucketListHandleFunc(w http.ResponseWriter, r *http.Request) {
	buckets := database.GetAllStatusBuckets()
	json.NewEncoder(w).Encode(buckets)
}

func StatusUserDateHTML(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	date := params["date"]

	u, err := database.GetUser(user)
	if err != nil {
		jww.ERROR.Printf("User: %s, does not exist.", user)
		http.Error(w, fmt.Errorf("User: %s, does not exist.", user).Error(), http.StatusNotFound)
		return
	}

	st := database.GetStatus(u, date)
	if len(st) == 0 {
		jww.ERROR.Printf("User: %s has no Status for %s", user, date)
		http.Error(w, fmt.Errorf("User: %s has no Status for %s", user, date).Error(), http.StatusNotFound)
		return
	}

	type Model struct {
		Date  string
		Users map[string][]database.StatusS
	}

	p := &Model{Date: date, Users: make(map[string][]database.StatusS)}
	p.Users[user] = st

	err = renderTemplate(w, "status.tmpl", p)
	if err != nil {
		jww.ERROR.Printf("Temp error: %+v", err)
	}
}

func StatusDateHTML(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	date := params["date"]

	st, err := database.GetAllStatusByDate(date)
	if err != nil {
		jww.ERROR.Printf("Cannot get all status' for %s", date)
		http.Error(w, fmt.Errorf("Cannot get all status' for %s", date).Error(), http.StatusNotFound)
		return
	}

	type Model struct {
		Date  string
		Users map[string][]database.StatusS
	}

	p := &Model{Date: date, Users: make(map[string][]database.StatusS)}

	for _, v := range st {
		if _, ok := p.Users[v.User.Nick]; !ok {
			p.Users[v.User.Nick] = []database.StatusS{v}
		} else {
			p.Users[v.User.Nick] = append(p.Users[v.User.Nick], v)
		}
	}

	err = renderTemplate(w, "status.tmpl", p)
	if err != nil {
		jww.ERROR.Printf("Temp error: %+v", err)
	}
}

func StatusHTML(w http.ResponseWriter, r *http.Request) {
	buckets := database.GetAllStatusBuckets()

	type Model struct {
		Date  string
		Dates []string
	}

	p := &Model{Date: database.DStamp, Dates: buckets}

	err := renderTemplate(w, "list.tmpl", p)
	if err != nil {
		jww.ERROR.Printf("Temp error: %+v", err)
	}

}

type StatusAPI struct {
	engine.API
}

func GetApi() *StatusAPI {
	return &StatusAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/status/":                      []engine.HandlersS{engine.MakeHandler("GET", StatusBucketListHandleFunc)},
				"/status/{user}":                []engine.HandlersS{engine.MakeHandler("PUT", StatusUserAddHandleFunc, middleware.AuthAndSelfOrAdmin...)},
				"/status/{date}":                []engine.HandlersS{engine.MakeHandler("GET", StatusAllGetHandleFunc)},
				"/status/{user}/{date}":         []engine.HandlersS{engine.MakeHandler("GET", StatusUserGetHandleFunc)},
				"/status/{user}/{date}/{index}": []engine.HandlersS{engine.MakeHandler("DELETE", StatusUserDelHandleFunc, middleware.AuthAndSelfOrAdmin...)},
				"/":              []engine.HandlersS{engine.MakeHandler("GET", StatusHTML)},
				"/{date}":        []engine.HandlersS{engine.MakeHandler("GET", StatusDateHTML)},
				"/{user}/{date}": []engine.HandlersS{engine.MakeHandler("GET", StatusUserDateHTML)},
			},
		},
	}
}
