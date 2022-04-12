package ldapx

import (
	"github.com/go-ldap/ldap/v3"
)

func (e *Entry) ReplaceAttributeValues(attr string, value []string) {
	v := e.Attributes.Get(attr)
	if v != nil && sameStringSlice(value, v.Values) {
		return
	}
	e.Attributes.PutEntryAttribute(ldap.NewEntryAttribute(attr, value))
	e.AddAttributeChange("replace", attr, value)
}

func (e *Entry) ReplaceAttributeValue(attr string, value string) {
	e.ReplaceAttributeValues(attr, []string{value})
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
