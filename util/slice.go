package util

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MapStringToSlice(m map[string]bool) []string {
	var l []string
	for v := range m {
		l = append(l, v)
	}
	return l
}

func ContinuousNums(start, count int, order string) (r []int) {
	if count <= 0 {
		return
	}
	for i := 0; i < count; i++ {
		if order == "desc" {
			if start-i < 0 {
				break
			}
			r = append(r, start-i)
		} else {
			r = append(r, start+i)
		}
	}
	return
}
