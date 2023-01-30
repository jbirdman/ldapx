package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

// ReplaceAttributeValuesIgnoreCase	removes all values from the attribute and adds the given values, ignoring case if ignoreCase is true.
func (e *Entry) ReplaceAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool) {
	v := e.Attributes.Get(attr)
	if v != nil && sameStringSlice(value, v.Values, ignoreCase) {
		return
	}
	e.Attributes.PutEntryAttribute(ldap.NewEntryAttribute(attr, value))
	e.AddAttributeChange("replace", attr, value)
}

// ReplaceAttributeValues removes all values from the attribute and adds the given values.
func (e *Entry) ReplaceAttributeValues(attr string, value []string) {
	e.ReplaceAttributeValuesIgnoreCase(attr, value, false)
}

// ReplaceAttributeValueIgnoreCase removes all values from the attribute and adds the given value, ignoring case if ignoreCase is true.
func (e *Entry) ReplaceAttributeValueIgnoreCase(attr string, value string, ignoreCase bool) {
	e.ReplaceAttributeValuesIgnoreCase(attr, []string{value}, ignoreCase)
}

func (e *Entry) ReplaceAttributeValue(attr string, value string) {
	e.ReplaceAttributeValueIgnoreCase(attr, value, false)
}

// sameStringSlice compares two string slices and returns true if they are the same, ignoring case if ignoreCase is true.
func sameStringSlice(x, y []string, ignoreCase bool) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		if ignoreCase {
			_x = strings.ToLower(_x)
		}
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		if ignoreCase {
			_y = strings.ToLower(_y)
		}
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
