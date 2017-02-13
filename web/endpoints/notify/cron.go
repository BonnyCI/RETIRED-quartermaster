package notify

import (
	"time"

	"github.com/robfig/cron"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/database"
)

var zone = map[string]int{
	"AEST": 10,
	"EST":  -5,
	"CST":  -6,
	"MST":  -7,
	"PST":  -8,
}

var StatusCron = "0 0 9 * * 1-5"

type TZJob struct {
	Group    database.GroupS
	TimeZone string
	Cron     string
}

func NewTzJob(g database.GroupS, c string) *TZJob {
	return &TZJob{
		Group:    g,
		TimeZone: g.Name,
		Cron:     c,
	}
}

func (t TZJob) Run() {
	doNotify(t.Group)
}

type Cron struct {
	Cron map[string]*cron.Cron
}

func (c Cron) Start() {
	for _, v := range c.Cron {
		go v.Start()
	}
}

func (c Cron) Stop() {
	for _, v := range c.Cron {
		v.Stop()
	}
}

func (c Cron) AddJob(j *TZJob) {
	c.Cron[j.TimeZone].AddJob(j.Cron, j)
}

func (c Cron) AddLocation(t string) {
	jww.DEBUG.Printf("Adding Cron for %s", t)
	loc := time.FixedZone(t, zone[t]*60*60)
	c.Cron[t] = cron.NewWithLocation(loc)
}

func BuildStatusCron(gs []database.GroupS) *Cron {
	c := &Cron{Cron: make(map[string]*cron.Cron)}

	for _, g := range gs {
		c.AddLocation(g.Name)
		c.AddJob(NewTzJob(g, StatusCron))
	}

	return c
}
