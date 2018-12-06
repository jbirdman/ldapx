package ldapx

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/ldap.v2"
)

func AddValue(attr string, values []string, value string) *ldap.PartialAttribute {
	return AddValues(attr, values, []string{value})
}

func AddValues(attr string, values []string, addValues []string) *ldap.PartialAttribute {
	currentValues := mapset.NewSetFromSlice(SliceToInterface(values))
	newValues := mapset.NewSetFromSlice(SliceToInterface(addValues))

	changes := newValues.Difference(currentValues)

	if changes.Cardinality() == 0 {
		return nil
	}

	return &ldap.PartialAttribute{Type: attr, Vals: InterfaceToSlice(changes.ToSlice())}
}

func RemoveValue(attr string, values []string, value string) *ldap.PartialAttribute {
	return RemoveValues(attr, values, []string{value})
}

func RemoveValues(attr string, values []string, removeValues []string) *ldap.PartialAttribute {
	currentValues := mapset.NewSetFromSlice(SliceToInterface(values))
	changes := mapset.NewSetFromSlice(SliceToInterface(removeValues))

	removes := currentValues.Difference(changes)

	if removes.Cardinality() == 0 {
		return nil
	}

	return &ldap.PartialAttribute{Type: attr, Vals: InterfaceToSlice(removes.ToSlice())}
}

