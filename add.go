package ldapx

import (
	"github.com/go-ldap/ldap/v3"
)

func (e *Entry) AddAttributeValues(attr string, value []string) {
	if len(value) == 0 {
		return
	}

	a, ok := e.Attributes[attr]
	if !ok {
		e.Attributes[attr] = ldap.NewEntryAttribute(attr, value)
		e.AddAttributeChange("add", attr, value)
		return
	}

	var addedValues []string // nolint:prealloc

	for _, o := range value {
		var found bool
		for _, d := range a.Values {
			if o == d {
				found = true
				break
			}
		}

		if found {
			continue
		}

		addedValues = append(addedValues, o)
		a.Values = append(a.Values, o)
	}

	if len(addedValues) > 0 {
		e.AddAttributeChange("add", attr, addedValues)
	}
}

func (e *Entry) AddAttributeValue(attr string, value string) {
	e.AddAttributeValues(attr, []string{value})
}
