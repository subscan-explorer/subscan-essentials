package util

import (
	"github.com/huandu/xstrings"

	"strings"
)

// Camel String
func CamelString(s string) string {
	return xstrings.ToPascalCase(s)
}

func UpperCamel(s string) string {
	if len(s) == 0 {
		return ""
	}
	s = strings.ToUpper(string(s[0])) + s[1:]
	return s
}

// String
func StringsExclude(a []string, b []string) []string {
	var refresh []string
	for _, v := range a {
		if !StringInSlice(v, b) {
			refresh = append(refresh, v)
		}
	}

	return refresh
}

func StringsIntersection(a []string, b []string) []string {
	var refresh []string
	for _, v := range a {
		if StringInSlice(v, b) {
			refresh = append(refresh, v)
		}
	}
	return refresh
}

func IfEmptyElse(a, b string) string {
	if strings.TrimSpace(a) == "" {
		return b
	}
	return a
}
