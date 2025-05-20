package util

func EnumPickOne(m map[string]string) string {
	for _, v := range m {
		if v != "" {
			return v
		}
	}
	return ""
}

func EnumPickOneInt(m map[string]int) int {
	for _, v := range m {
		if v != 0 {
			return v
		}
	}
	return 0
}

func EnumStringKey(m map[string]string) string {
	for k := range m {
		return k
	}
	return ""
}

func EnumKey(m map[string]interface{}) string {
	for k := range m {
		return k
	}
	return ""
}
