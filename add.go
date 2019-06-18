package ldapx

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/ldap.v2"
)

func (e *Entry) AddAttributeValues(attr string, value []string) {
	var values []string

	a, ok := e.Attributes[attr]
	if !ok {
		a = ldap.NewEntryAttribute(attr, value)
		values = value
	} else {
		oldValues := mapset.NewSetFromSlice(SliceToInterface(a.Values))
		newValues := mapset.NewSetFromSlice(SliceToInterface(value))

		valueSet := newValues.Difference(oldValues)

		for _, v := range valueSet.ToSlice() {
			a.Values = append(a.Values, v.(string))
			values = append(values, v.(string))
		}
	}

	if len(values) > 0 {
		e.Attributes[attr] = a
		e.AddAttributeChange("add", attr, values)
	}
}

func (e *Entry) AddAttributeValue(attr string, value string) {
	e.AddAttributeValues(attr, []string{value})
}
