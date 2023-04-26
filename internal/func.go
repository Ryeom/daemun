package internal

func Contains(l []string, p string) bool {
	for _, v := range l {
		if v == p {
			return true
		}
	}
	return false
}

func UniqueList(s []string) []string {
	keys := make(map[string]struct{})
	res := make([]string, 0)
	for _, val := range s {
		if _, ok := keys[val]; ok {
			continue
		} else {
			keys[val] = struct{}{}
			res = append(res, val)
		}
	}
	return res
}
