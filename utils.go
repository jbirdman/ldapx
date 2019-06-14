package ldapx

import "github.com/deckarep/golang-set"

func SliceToInterface(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = v
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

func ContainsAny(a []string, b []string) bool {
	s1 := mapset.NewSetFromSlice(SliceToInterface(a))
	s2 := mapset.NewSetFromSlice(SliceToInterface(b))
	return s2.IsSubset(s1)
}
