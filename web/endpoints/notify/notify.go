package notify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/bot"
	"github.com/bonnyci/quartermaster/database"
	"github.com/bonnyci/quartermaster/web/engine"
)

func doNotify(g database.GroupS) {
	jww.INFO.Printf("Notiifying members of %s, %v", g.Name, g.Members)
	i := bot.GetIrc()

	for _, v := range g.Members {
		i.Sendf(v.Nick, "Good morning %s, your status report time is now open. Usage: `status add <status message>`", v.Nick)
	}
}

func NotifyGroupHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	g, err := database.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}

	doNotify(g)
	json.NewEncoder(w).Encode(g)
}

type NotifyAPI struct {
	engine.API
}

func GetApi() *NotifyAPI {
	return &NotifyAPI{
		engine.APIBase{
			Handlers: engine.HandlersT{
				"/notify/{group}": []engine.HandlersS{engine.MakeHandler("GET", NotifyGroupHandleFunc)},
			},
		},
	}
}

func NotifyCron() *Cron {
	gs, _ := database.ListGroups()
	g, _ := database.GetGroup("Admin")

	gs = database.RemoveGroup(gs, g)
	jww.DEBUG.Printf("%+v", gs)
	return BuildStatusCron(gs)
}
