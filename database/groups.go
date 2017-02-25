package database

import (
	"fmt"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type GroupS struct {
	ID      int     `json:"id"`
	Name    string  `json:"name" storm:"index,unique"`
	Members []UserS `json:"members,omitempty"`
}

func (s GroupS) String() string {
	var ms []string

	for _, v := range s.Members {
		ms = append(ms, v.Nick)
	}

	m := strings.Join(ms, ",")

	return "(" + s.Name + " M:" + m + ")"
}

func (s GroupS) Compare(d DataS) bool {
	c := d.(GroupS)

	if s.Name != c.Name {
		return false
	}

	if s.Members == nil && c.Members == nil {
		return true
	}

	if s.Members == nil || c.Members == nil {
		return false
	}

	if len(s.Members) != len(c.Members) {
		return false
	}

	for i := range s.Members {
		if !s.Members[i].Compare(c.Members[i]) {
			return false
		}
	}

	return true
}

func (s GroupS) Save() error {
	db := GetInstance()
	jww.DEBUG.Printf("Saving: %+v", s)
	if err := db.DbObj.Save(&s); err != nil {
		jww.ERROR.Printf("Failure to Save: %+v", s)
		return err
	}
	return nil
}

func (s GroupS) Delete() error {
	db := GetInstance()
	jww.DEBUG.Printf("Deleting: %+v", s)
	if err := db.DbObj.DeleteStruct(&s); err != nil {
		jww.ERROR.Printf("Failure to Delete: %+v", s)
		return err
	}
	return nil
}

func (s GroupS) Update() error {
	db := GetInstance()
	jww.DEBUG.Printf("Updating: %+v", s)
	if err := db.DbObj.Update(&s); err != nil {
		jww.ERROR.Printf("Failure to Update: %+v", s)
		return err
	}
	return nil
}

func ListGroups() ([]GroupS, error) {
	var dg []GroupS
	if err := GetAll(&dg); err != nil {
		return dg, err
	}
	return dg, nil
}

func GetGroup(n string) (GroupS, error) {
	var g GroupS
	if err := One("Name", n, &g); err != nil {
		return g, fmt.Errorf("Group: %s does not exist.", n)
	}
	return g, nil
}

func addGroup(gs string) error {
	g := GroupS{Name: gs}
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
	var g GroupS
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

func ModifyGroups(gs []GroupS) error {
	for _, g := range gs {
		if err := g.Update(); err != nil {
			return err
		}
	}
	return nil
}

func AddUsersToGroups(gs []string, us []string) error {
	var err error
	var grs []GroupS
	for _, g := range gs {
		var gr GroupS
		if gr, err = GetGroup(g); err != nil {
			return err
		}
		for _, v := range us {
			var u UserS
			if u, err = GetUser(v); err != nil {
				return fmt.Errorf("Attempting to add %s to %s, user does not exist.", v, g)
			}
			if UserInGroup(gr, u) {
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
	var grs []GroupS
	for _, g := range gs {
		var gr GroupS
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

func UserInGroup(g GroupS, u UserS) bool {
	for _, v := range g.Members {
		if v.Nick == u.Nick {
			return true
		}
	}
	return false
}

func RemoveGroup(list []GroupS, rm GroupS) []GroupS {
	for k, v := range list {
		if ok := v.Compare(rm); ok {
			return append(list[:k], list[k+1:]...)
		}
	}
	return list
}
