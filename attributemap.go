package ldapx

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// AttributeMap is a map of attributes.
type AttributeMap map[string]*ldap.EntryAttribute

// NewAttributeMap creates a new AttributeMap.
func NewAttributeMap() AttributeMap {
	return make(AttributeMap)
}

// AttributeNames returns the names of the attributes in the map.
func (m AttributeMap) AttributeNames() []string {
	names := make([]string, 0, len(m))
	for _, a := range m {
		names = append(names, a.Name)
	}
	return names
}

// Get returns the attribute with the given name.
func (m AttributeMap) Get(attr string) *ldap.EntryAttribute {
	v := m[strings.ToLower(attr)]
	return v
}

// Put adds the attribute to the map.
func (m AttributeMap) Put(attr string, value *ldap.EntryAttribute) {
	m[strings.ToLower(attr)] = value
}

// PutEntryAttribute adds the attribute from the given LDAP entry to the map
func (m AttributeMap) PutEntryAttribute(value *ldap.EntryAttribute) {
	m[strings.ToLower(value.Name)] = value
}

// Delete removes the attribute with the given name from the map.
func (m AttributeMap) Delete(attr string) {
	delete(m, strings.ToLower(attr))
}

// Rename renames the attribute with the given name to the new name.
func (m AttributeMap) Rename(from, to string) {
	// If the names are the same, do nothing
	if strings.EqualFold(from, to) {
		return
	}

	// Get the attribute
	v := m[strings.ToLower(from)]

	// If it doesn't exist, do nothing
	if v == nil {
		return
	}

	// Rename the attribute in the value
	v.Name = to

	// Add "new" version
	m.PutEntryAttribute(v)

	// And delete old attribute
	m.Delete(from)
}

func (m AttributeMap) AttributeExists(attr string) bool {
	_, ok := m[strings.ToLower(attr)]
	return ok
}
