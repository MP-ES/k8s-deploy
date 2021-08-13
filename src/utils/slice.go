package utils

import (
	"strings"
)

func SliceMapStrFunction(vs []string, f func(string) string) []string {
	if vs == nil || f == nil {
		return vs
	}

	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func SliceRemoveEmptyElements(s []string) []string {
	if s == nil {
		return s
	}

	var r []string = []string{}
	for _, str := range s {
		if strings.TrimSpace(str) != "" && str != "null" {
			r = append(r, str)
		}
	}
	return r
}

func SliceRemoveDuplicateElements(s []string) []string {
	if s == nil {
		return s
	}

	keys := make(map[string]bool)
	list := []string{}
	for _, e := range s {
		if _, value := keys[e]; !value {
			keys[e] = true
			list = append(list, e)
		}
	}
	return list
}
