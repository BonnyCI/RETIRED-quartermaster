package lib

import (
	"fmt"

	"github.com/bonnyci/quartermaster/database"
)

func ListGroups() ([]database.GroupS, error) {
	var dg []database.GroupS
	if err := database.GetAll(&dg); err != nil {
		return dg, err
	}
	return dg, nil
}

func GetGroup(n string) (database.GroupS, error) {
	var g database.GroupS
	if err := database.One("Name", n, &g); err != nil {
		return g, fmt.Errorf("Group: %s does not exist.", n)
	}
	return g, nil
}

func addGroup(gs string) error {
	g := database.GroupS{Name: gs}
	if err := g.Save(); err != nil {
		return err
	}
	return nil
}

func AddGroups(gs []string) error {
	for _, g := range gs {
		if err := addGroup(g); err != nil {
			return err
		}
	}
	return nil
}

func delGroup(gs string) error {
	var g database.GroupS
	var err error
	if g, err = GetGroup(gs); err != nil {
		return err
	}
	if err = g.Delete(); err != nil {
		return err
	}
	return nil
}

func DelGroups(gs []string) error {
	for _, g := range gs {
		if err := delGroup(g); err != nil {
			return err
		}
	}
	return nil
}

func ModifyGroups(gs []database.GroupS) error {
	for _, g := range gs {
		if err := g.Update(); err != nil {
			return err
		}
	}
	return nil
}

func AddUsersToGroups(gs []string, us []string) error {
	var err error
	var grs []database.GroupS
	for _, g := range gs {
		var gr database.GroupS
		if gr, err = GetGroup(g); err != nil {
			return err
		}
		for _, v := range us {
			var u database.UserS
			if u, err = GetUser(v); err != nil {
				return fmt.Errorf("Attempting to add %s to %s, user does not exist.", v, g)
			}
			if UserInGroup(g, u) {
				continue
			}
			gr.Members = append(gr.Members, u)
		}
		grs = append(grs, gr)
	}
	ModifyGroups(grs)
	return nil

}

func DelUsersFromGroups(gs []string, us []string) error {
	var err error
	var grs []database.GroupS
	for _, g := range gs {
		var gr database.GroupS
		if gr, err = GetGroup(g); err != nil {
			return err
		}
		for _, v := range us {
			u, _ := GetUser(v)
			gr.Members = RemoveUser(gr.Members, u)
		}
		grs = append(grs, gr)
	}
	ModifyGroups(grs)
	return nil

}

func UserInGroup(n string, u database.UserS) bool {
	var err error
	var g database.GroupS
	if g, err = GetGroup(n); err != nil {
		return false
	}
	for _, v := range g.Members {
		if v.Nick == u.Nick {
			return true
		}
	}
	return false
}
