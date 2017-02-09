package notify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/bot"
	"github.com/bonnyci/quartermaster/lib"
	"github.com/bonnyci/quartermaster/web/engine"
)

func NotifyGroupHandleFunc(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	group := params["group"]

	g, err := lib.GetGroup(group)
	if err != nil {
		jww.ERROR.Printf("Group: %s, does not exist.", group)
		http.Error(w, fmt.Errorf("Group: %s, does not exist.", group).Error(), http.StatusNotFound)
		return
	}

	i := bot.GetIrc()

	for _, v := range g.Members {
		i.Sendf(v.Nick, "Good morning %s, your status report time is now open.", v.Nick)
	}

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
