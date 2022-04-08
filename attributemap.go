package ldapx

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type AttributeMap struct {
	Attrs map[string]*ldap.EntryAttribute `json:"Attrs"`
}

func NewAttributeMap() AttributeMap {
	return AttributeMap{
		make(map[string]*ldap.EntryAttribute),
	}
}

func (m AttributeMap) AttributeNames() []string {
	var names []string
	for _, a := range m.Attrs {
		names = append(names, a.Name)
	}
	return names
}

func (m AttributeMap) Get(attr string) *ldap.EntryAttribute {
	v, _ := m.Attrs[strings.ToLower(attr)]
	return v
}

func (m *AttributeMap) Put(attr string, value *ldap.EntryAttribute) {
	m.Attrs[strings.ToLower(attr)] = value
}

func (m *AttributeMap) PutEntryAttribute(value *ldap.EntryAttribute) {
	m.Put(value.Name, value)
}

func (m *AttributeMap) Delete(attr string) {
	delete(m.Attrs, strings.ToLower(attr))
}

func (m *AttributeMap) Rename(from, to string) {
	if strings.ToLower(from) == strings.ToLower(to) {
		return
	}

	v := m.Get(from)

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
