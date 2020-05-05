package ldapx

import mapset "github.com/deckarep/golang-set"

func (e *Entry) SyncAttributeValues(attr string, values []string) {
	current := mapset.NewSetFromSlice(stringsToInterfaceSlice(e.GetAttributeValues(attr)))
	newset := mapset.NewSetFromSlice(stringsToInterfaceSlice(values))

	e.AddAttributeValues(attr, interfaceSliceToStrings(newset.Difference(current).ToSlice()))
	e.DeleteAttributeValues(attr, interfaceSliceToStrings(current.Difference(newset).ToSlice()))
}

func stringsToInterfaceSlice(items []string) []interface{} {
	var intItems = make([]interface{}, len(items))

	for i, v := range items {
		intItems[i] = v
	}

	return intItems
}

func interfaceSliceToStrings(items []interface{}) []string {
	var stringItems = make([]string, len(items))

	for i, v := range items {
		stringItems[i] = v.(string)
	}

	return stringItems
}
