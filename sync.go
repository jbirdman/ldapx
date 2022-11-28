package ldapx

import mapset "github.com/deckarep/golang-set"

func (e *Entry) SyncAttributeValues(attr string, values []string) {
	e.SyncAttributeValuesIgnoreCase(attr, values, false)
}

func (e *Entry) SyncAttributeValuesIgnoreCase(attr string, values []string, ignoreCase bool) {
	current := mapset.NewSetFromSlice(SliceToInterface(e.GetAttributeValues(attr)))
	newset := mapset.NewSetFromSlice(SliceToInterface(values))

	e.AddAttributeValues(attr, InterfaceToSlice(newset.Difference(current).ToSlice()))
	e.DeleteAttributeValues(attr, InterfaceToSlice(current.Difference(newset).ToSlice()))
}
