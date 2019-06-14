package ldapx

import (
	"gopkg.in/ldap.v2"
	"reflect"
)

func (e *Entry) ReplaceAttributeValues(attr string, value []string) {
	if reflect.DeepEqual(value, e.Attributes[attr].Values) {
		return
	}
	e.Attributes[attr] = ldap.NewEntryAttribute(attr, value)
	e.AddAttributeChange("replace", attr, value)
}

func (e *Entry) ReplaceAttributeValue(attr string, value string) {
	e.ReplaceAttributeValues(attr, []string{value})
}

