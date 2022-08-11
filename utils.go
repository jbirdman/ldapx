package ldapx

import (
	"strings"

	"github.com/deckarep/golang-set"
)

func SliceToInterface(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = v
	}

	return s
}

func SliceToInterfaceFold(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = strings.ToLower(v)
	}

	return s
}

func InterfaceToSlice(t []interface{}) []string {
	s := make([]string, len(t))

	for i, v := range t {
		s[i] = v.(string)
	}

	return s
}

func InterfaceToSliceFold(t []interface{}) []string {
	s := make([]string, len(t))

	for i, v := range t {
		s[i] = strings.ToLower(v.(string))
	}
	return s
}

func StringSliceFold(t []string) []string {
	return MapStringSlice(t, strings.ToLower)
}

func MapStringSlice(t []string, f func(string) string) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = f(v)
	}
	return s
}

func ContainsAny(a []string, b []string) bool {
	s1 := mapset.NewSetFromSlice(SliceToInterface(a))
	s2 := mapset.NewSetFromSlice(SliceToInterface(b))
	return s2.IsSubset(s1)
}

func ContainsAnyFold(a []string, b []string) bool {
	return ContainsAny(MapStringSlice(a, strings.ToLower), MapStringSlice(b, strings.ToLower))
}
