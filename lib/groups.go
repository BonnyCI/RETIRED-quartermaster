package lib

import (
	"github.com/pschwartz/quartermaster/database"
)

func ListGroups() []database.GroupS {
	var dg []database.GroupS
	database.GetAll(&dg)
	return dg
}

func GetGroup(n string) database.GroupS {
	var g database.GroupS
	database.One("Name", n, &g)
	return g
}

func addGroup(gs string) {
	g := database.GroupS{Name: gs}
	g.Save()
}

func AddGroups(gs []string) {
	for _, g := range gs {
		addGroup(g)
	}
}

func delGroup(gs string) {
	g := GetGroup(gs)
	g.Delete()
}

func DelGroups(gs []string) {
	for _, g := range gs {
		delGroup(g)
	}
}

func ModifyGroups(gs []database.GroupS) {
	for _, g := range gs {
		g.Update()
	}
}

func AddUsersToGroups(t string, gs []string, us []string) {
	var grs []database.GroupS
	for _, g := range gs {
		gr := GetGroup(g)
		switch t {
		case "Member":
			gr.Members = append(gr.Members, us...)
		case "Admin":
			gr.Admins = append(gr.Admins, us...)
		}
		grs = append(grs, gr)
	}
	ModifyGroups(grs)

}

func DelUsersFromGroups(t string, gs []string, us []string) {
	var grs []database.GroupS
	for _, g := range gs {
		gr := GetGroup(g)
		for _, u := range us {
			switch t {
			case "Member":
				gr.Members = Remove(gr.Members, u)
			case "Admin":
				gr.Admins = Remove(gr.Admins, u)
			}
		}
		grs = append(grs, gr)
	}
	ModifyGroups(grs)

}
