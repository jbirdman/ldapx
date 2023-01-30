package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

// AddAttributeValuesIgnoreCase adds the given values to the attribute.
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

		// If the value is found, skip it
		if found {
			continue
		}

		// If the value was not found, add it to the added values
		addedValues = append(addedValues, o)
		a.Values = append(a.Values, o)
	}

	// Add the attribute change to the entry
	if len(addedValues) > 0 {
		e.AddAttributeChange("add", attr, addedValues)
	}
}

// AddAttributeValues adds the given values to the attribute.
func (e *Entry) AddAttributeValues(attr string, value []string) {
	e.AddAttributeValuesIgnoreCase(attr, value, false)
}

// AddAttributeValueIgnoreCase adds the given value to the attribute.
func (e *Entry) AddAttributeValueIgnoreCase(attr string, value string, ignoreCase bool) {
	e.AddAttributeValuesIgnoreCase(attr, []string{value}, ignoreCase)
}

// AddAttributeValue adds the given value to the attribute.
func (e *Entry) AddAttributeValue(attr string, value string) {
	e.AddAttributeValueIgnoreCase(attr, value, false)
}
