package ldapx

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/jbirdman/caseinsensitiveset"
)

// SyncAttributeValues will add and remove values from the attribute to match the
// values provided.
func (e *Entry) SyncAttributeValues(attr string, values []string) {
	e.SyncAttributeValuesIgnoreCase(attr, values, false)
}

// SyncAttributeValuesIgnoreCase will add and remove values from the attribute to
// match the values provided, ignoring case if ignoreCase is true.
func (e *Entry) SyncAttributeValuesIgnoreCase(attr string, values []string, ignoreCase bool) {
	if ignoreCase {
		e.syncAttributeValuesIgnoreCase(attr, values)
	} else {
		e.syncAttributeValuesCaseSensitive(attr, values)
	}
}

// syncAttributeValuesCaseSensitive will add and remove values from the attribute
// to match the values provided, ignoring case.
func (e *Entry) syncAttributeValuesIgnoreCase(attr string, value []string) {
	currentValues := caseinsensitiveset.NewCaseInsensitiveSet(e.GetAttributeValues(attr)...)
	newValues := caseinsensitiveset.NewCaseInsensitiveSet(value...)

	e.AddAttributeValues(attr, newValues.Difference(currentValues).ToSlice())
	e.DeleteAttributeValues(attr, currentValues.Difference(newValues).ToSlice())
}

// syncAttributeValuesCaseSensitive will add and remove values from the attribute
// to match the values provided.
func (e *Entry) syncAttributeValuesCaseSensitive(attr string, value []string) {
	currentValues := mapset.NewSet(e.GetAttributeValues(attr)...)
	newValues := mapset.NewSet(value...)

	e.AddAttributeValues(attr, newValues.Difference(currentValues).ToSlice())
	e.DeleteAttributeValues(attr, currentValues.Difference(newValues).ToSlice())
}
