package endpoints

import (
	"github.com/bonnyci/quartermaster/web/endpoints/backup"
	"github.com/bonnyci/quartermaster/web/endpoints/groups"
	"github.com/bonnyci/quartermaster/web/endpoints/notify"
	"github.com/bonnyci/quartermaster/web/endpoints/status"
	"github.com/bonnyci/quartermaster/web/endpoints/users"
	"github.com/bonnyci/quartermaster/web/engine"
)

func BackupHTTP() {
	be := engine.GetEngine("backup")
	be.SetAddr(":12800")
	backup.GetApi().AddToEngine(be)
	be.Start()
}

func ApiHTTP() {
	be := engine.GetEngine("api")
	be.SetAddr(":8888")
	groups.GetApi().AddToEngine(be)
	notify.GetApi().AddToEngine(be)
	users.GetApi().AddToEngine(be)
	status.GetApi().AddToEngine(be)
	be.Start()
}
