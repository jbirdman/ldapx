package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type AttributeMap map[string]*ldap.EntryAttribute

func NewAttributeMap() AttributeMap {
	return make(AttributeMap)
}

func (m AttributeMap) AttributeNames() []string {
	var names = make([]string, 0, len(m))
	for _, a := range m {
		names = append(names, a.Name)
	}
	return names
}

func (m AttributeMap) Get(attr string) *ldap.EntryAttribute {
	v := m[strings.ToLower(attr)]
	return v
}

func (m AttributeMap) Put(attr string, value *ldap.EntryAttribute) {
	m[strings.ToLower(attr)] = value
}

func (m AttributeMap) PutEntryAttribute(value *ldap.EntryAttribute) {
	m[strings.ToLower(value.Name)] = value
}

func (m AttributeMap) Delete(attr string) {
	delete(m, strings.ToLower(attr))
}

func (m AttributeMap) Rename(from, to string) {
	if strings.EqualFold(from, to) {
		return
	}

	v := m[strings.ToLower(from)]

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
