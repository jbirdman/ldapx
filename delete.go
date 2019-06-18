package ldapx

import (
	"gopkg.in/ldap.v2"
)

func (e *Entry) DeleteAttributeValues(attr string, value []string) {
	a, ok := e.Attributes[attr]
	if !ok {
		return
	}
	var deletedValues []string
	var remainingValues []string

	for _, o := range a.Values {
		var found bool
		for _, d := range value {
			if o == d {
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
		delete(e.Attributes, attr)
		e.AddAttributeChange("delete", attr, nil)
	} else {
		e.Attributes[attr] = ldap.NewEntryAttribute(attr, remainingValues)
		e.AddAttributeChange("delete", attr, deletedValues)
	}
}

func (e *Entry) DeleteAttributeValue(attr string, value string) {
	e.DeleteAttributeValues(attr, []string{value})
}

func (e *Entry) DeleteAttribute(attr string) {
	delete(e.Attributes, attr)
	e.AddAttributeChange("delete", attr, nil)
}
