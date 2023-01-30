package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

// DeleteAttributeValuesIgnoreCase deletes the given values from the attribute.
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

		// If the value was found, add it to the deleted values, otherwise add it to the remaining values
		if found {
			deletedValues = append(deletedValues, o)
		} else {
			remainingValues = append(remainingValues, o)
		}
	}

	// Nothing to do
	if len(deletedValues) == 0 {
		return
	}

	// All values were deleted
	if len(remainingValues) == 0 {
		e.Attributes.Delete(attr)
		e.AddAttributeChange("delete", attr, nil)
	} else {
		// Some values were deleted
		e.Attributes.PutEntryAttribute(ldap.NewEntryAttribute(attr, remainingValues))
		e.AddAttributeChange("delete", attr, deletedValues)
	}
}

// DeleteAttributeValues deletes the given values from the attribute.
func (e *Entry) DeleteAttributeValues(attr string, value []string) {
	e.DeleteAttributeValuesIgnoreCase(attr, value, false)
}

// DeleteAttributeValueIgnoreCase deletes the given value from the attribute.
func (e *Entry) DeleteAttributeValueIgnoreCase(attr string, value string, ignoreCase bool) {
	e.DeleteAttributeValuesIgnoreCase(attr, []string{value}, ignoreCase)
}

// DeleteAttributeValue deletes the given value from the attribute.
func (e *Entry) DeleteAttributeValue(attr string, value string) {
	e.DeleteAttributeValueIgnoreCase(attr, value, false)
}

// DeleteAttribute deletes the attribute.
func (e *Entry) DeleteAttribute(attr string) {
	e.Attributes.Delete(attr)
	e.AddAttributeChange("delete", attr, nil)
}
