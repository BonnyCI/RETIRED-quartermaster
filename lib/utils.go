package lib

import (
	"encoding/json"
	"io"
	"time"

	"github.com/bonnyci/quartermaster/database"
)

const DFMT = "2006-01-02"

var DStamp = time.Now().UTC().Format(DFMT)

func Remove(list []database.UserS, rm database.UserS) []database.UserS {
	for k, v := range list {
		if v == rm {
			return append(list[:k], list[k+1:]...)
		}
	}
	return list
}

func RemoveIndex(list []string, i int) []string {
	for k := range list {
		if k == i {
			return append(list[:k], list[k+1:]...)
		}
	}
	return list
}

func In(v string, a []string) (ok bool, i int) {
	for i = range a {
		if ok = a[i] == v; ok {
			return
		}
	}
	return
}

func Build(in io.Reader, out interface{}) error {
	if err := json.NewDecoder(in).Decode(out); err != nil {
		return err
	}
	return nil
}
