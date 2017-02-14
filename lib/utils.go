package lib

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
