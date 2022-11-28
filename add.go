package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func (e *Entry) AddAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool) {
	if len(value) == 0 {
		return
	}

	a := e.Attributes.Get(attr)
	if a == nil {
		e.Attributes.PutEntryAttribute(ldap.NewEntryAttribute(attr, value))
		e.AddAttributeChange("add", attr, value)
		return
	}

	var addedValues []string //nolint:prealloc

	for _, o := range value {
		var found bool
		for _, d := range a.Values {
			if (o == d) || (ignoreCase && strings.EqualFold(o, d)) {
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

func (e *Entry) AddAttributeValues(attr string, value []string) {
	e.AddAttributeValuesIgnoreCase(attr, value, false)
}

func (e *Entry) AddAttributeValueIgnoreCase(attr string, value string, ignoreCase bool) {
	e.AddAttributeValuesIgnoreCase(attr, []string{value}, ignoreCase)
}

func (e *Entry) AddAttributeValue(attr string, value string) {
	e.AddAttributeValueIgnoreCase(attr, value, false)
}
