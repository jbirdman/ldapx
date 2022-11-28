package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func (e *Entry) DeleteAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool) {
	a := e.Attributes.Get(attr)
	if a == nil {
		return
	}
	var deletedValues []string
	var remainingValues []string

	for _, o := range a.Values {
		var found bool
		for _, d := range value {
			if (o == d) || (ignoreCase && strings.EqualFold(o, d)) {
				found = true
				break
			}
		}

		if found {
			deletedValues = append(deletedValues, o)
		} else {
			remainingValues = append(remainingValues, o)
		}
	}

	if len(deletedValues) == 0 {
		return
	}

	if len(remainingValues) == 0 {
		e.Attributes.Delete(attr)
		e.AddAttributeChange("delete", attr, nil)
	} else {
		e.Attributes.PutEntryAttribute(ldap.NewEntryAttribute(attr, remainingValues))
		e.AddAttributeChange("delete", attr, deletedValues)
	}
}

func (e *Entry) DeleteAttributeValues(attr string, value []string) {
	e.DeleteAttributeValuesIgnoreCase(attr, value, false)
}

func (e *Entry) DeleteAttributeValueIgnoreCase(attr string, value string, ignoreCase bool) {
	e.DeleteAttributeValuesIgnoreCase(attr, []string{value}, ignoreCase)
}

func (e *Entry) DeleteAttributeValue(attr string, value string) {
	e.DeleteAttributeValueIgnoreCase(attr, value, false)
}

func (e *Entry) DeleteAttribute(attr string) {
	e.Attributes.Delete(attr)
	e.AddAttributeChange("delete", attr, nil)
}
