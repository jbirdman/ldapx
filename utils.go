package ldapx

import (
	"github.com/deckarep/golang-set/v2"
	"strings"
)

// SliceToInterface converts a string slice to an interface slice.
func SliceToInterface(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = v
	}

	return s
}

// SliceToInterfaceFold converts a string slice to an interface slice, lowercasing
func SliceToInterfaceFold(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = strings.ToLower(v)
	}

	return s
}

// InterfaceToSlice converts an interface slice to a string slice.
func InterfaceToSlice(t []interface{}) []string {
	s := make([]string, len(t))

	for i, v := range t {
		s[i] = v.(string)
	}

	return s
}

// InterfaceToSliceFold converts an interface slice to a string slice, lowercasing
func InterfaceToSliceFold(t []interface{}) []string {
	s := make([]string, len(t))

	for i, v := range t {
		s[i] = strings.ToLower(v.(string))
	}
	return s
}

// StringSliceFold converts a string slice to a lowercased string slice
func StringSliceFold(t []string) []string {
	return MapStringSlice(t, strings.ToLower)
}

// MapStringSlice applies a function to each element of a string slice
func MapStringSlice(t []string, f func(string) string) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = f(v)
	}
	return s
}

func ContainsAny(a []string, b []string) bool {
	s1 := mapset.NewSet(a...)
	s2 := mapset.NewSet(b...)
	return s2.IsSubset(s1)
}

func ContainsAnyFold(a []string, b []string) bool {
	return ContainsAny(MapStringSlice(a, strings.ToLower), MapStringSlice(b, strings.ToLower))
}
