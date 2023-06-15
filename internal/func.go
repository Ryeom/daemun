package internal

func Contains(l []any, p any) bool {
	for _, v := range l {
		if v == p {
			return true
		}
	}
	return false
}

func Unique(s []any) []any {
	keys := make(map[any]struct{})
	res := make([]any, 0)
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
