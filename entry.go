package ldapx

import (
	"fmt"
	"gopkg.in/ldap.v2"
)

const (
	CHANGE_ADD    = "add"
	CHANGE_UPDATE = "update"
	CHANGE_DELETE = "delete"
)

type Entry struct {
	ChangeType string
	DN         string
	Attributes map[string]*ldap.EntryAttribute
	Changes    []AttributeChange
}

type MutableEntry interface {
	AddAttribute(attr string, value []string)
	ReplaceAttribute(attr string, value []string)
	DeleteAttribute(attr string, value []string)
	Update(conn *Conn) error
}

type AttributeChange struct {
	Action string
	Attr   string
	Value  []string
}

func NewEntry(entry *ldap.Entry) *Entry {
	var dn string
	attributes := make(map[string]*ldap.EntryAttribute)

	var changeType = CHANGE_ADD
	if entry != nil {
		dn = entry.DN
		changeType = CHANGE_UPDATE

		// Copy in the attributes from the ldap entry
		for _, a := range entry.Attributes {
			attributes[a.Name] = a
		}
	}

	return &Entry{ChangeType: changeType, DN: dn, Attributes: attributes}
}

func (e *Entry) ToLdapEntry() *ldap.Entry {
	attrs := make(map[string][]string)

	for k, v := range e.Attributes {
		attrs[k] = v.Values
	}

	return ldap.NewEntry(e.DN, attrs)
}

func (e *Entry) Print() {
	fmt.Printf("DN: %s\n", e.DN)
	for _, attr := range e.Attributes {
		attr.Print()
	}
}

func (e *Entry) Clear() {
	e.ChangeType = CHANGE_ADD
	e.Changes = nil
}

func (e *Entry) GetAttributeValues(attribute string) []string {
	return e.Attributes[attribute].Values
}

func (e *Entry) GetAttributeValue(attribute string) string {
	values := e.GetAttributeValues(attribute)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (e *Entry) AddAttributeChange(action string, attr string, value []string) {
	e.Changes = append(e.Changes, AttributeChange{Action: action, Attr: attr, Value: value})
}
func (e *Entry) AddAttribute(attr string, value []string) {
	a, ok := e.Attributes[attr]
	if !ok {
		a = ldap.NewEntryAttribute(attr, value)
	} else {
		a.Values = append(a.Values, value...)
	}
	e.Attributes[attr] = a
	e.AddAttributeChange("add", attr, value)
}

func (e *Entry) ReplaceAttribute(attr string, value []string) {
	e.Attributes[attr] = ldap.NewEntryAttribute(attr, value)
	e.AddAttributeChange("replace", attr, value)
}

func (e *Entry) DeleteAttribute(attr string, value []string) {
	old, ok := e.Attributes[attr]
	if !ok {
		return
	}
	if len(value) > 0 {
		a := ldap.NewEntryAttribute(attr, nil)

		for _, v := range old.Values {
			found := false
			for _, nv := range value {
				if nv == v {
					found = true
					break
				}
			}

			if !found {
				a.Values = append(a.Values, v)
			}
		}
		e.Attributes[attr] = a
	} else {
		delete(e.Attributes, attr)
	}

	e.AddAttributeChange("delete", attr, value)
}

func (e Entry) Update(conn *Conn) error {
	if len(e.Changes) == 0 {
		return nil
	}

	switch e.ChangeType {
	case CHANGE_ADD:
		return conn.Add(buildAddRequest(e.DN, e.Changes))
	case CHANGE_UPDATE:
		return conn.Modify(buildModifyRequest(e.DN, e.Changes))
	case CHANGE_DELETE:
		return conn.Del(buildDelRequest(e.DN))
	}

	return nil
}

func buildAddRequest(dn string, changes []AttributeChange) *ldap.AddRequest {
	r := NewAddRequest(dn)

	for _, change := range changes {
		if change.Action != "add" && change.Action != "replace" {
			continue
		}

		r.Attribute(change.Attr, change.Value)
	}

	return r
}

func buildModifyRequest(dn string, changes []AttributeChange) *ldap.ModifyRequest {
	r := NewModifyRequest(dn)

	for _, change := range changes {
		switch change.Action {
		case "add":
			r.Add(change.Attr, change.Value)
		case "replace":
			r.Replace(change.Attr, change.Value)
		case "delete":
			r.Delete(change.Attr, change.Value)
		}
	}

	return r
}

func buildDelRequest(dn string) *ldap.DelRequest {
	return NewDelRequest(dn, nil)
}
