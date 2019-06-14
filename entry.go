package ldapx

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"gopkg.in/ldap.v2"
	"reflect"
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
	AddAttributeValue(attr string, value string)
	AddAttributeValues(attr string, value []string)
	ReplaceAttributeValue(attr string, value string)
	ReplaceAttributeValues(attr string, value []string)
	DeleteAttributeValue(attr string, value string)
	DeleteAttributeValues(attr string, value []string)
	Update(conn *Conn) error
}

type AttributeChange struct {
	Action string
	Attr   string
	Value  []string
}

func NewEntry(dn string) *Entry {
	return &Entry{DN: dn, Attributes: make(map[string]*ldap.EntryAttribute), ChangeType: CHANGE_ADD}
}

func NewEntryFromLdapEntry(entry *ldap.Entry) *Entry {
	e := NewEntry("")

	if entry != nil {
		e.DN = entry.DN
		e.ChangeType = CHANGE_UPDATE

		// Copy in the attributes from the ldap entry
		for _, a := range entry.Attributes {
			e.Attributes[a.Name] = a
		}
	}

	return e
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

func (e *Entry) ResetChanges() {
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

func (e *Entry) AddAttributeValues(attr string, value []string) {
	var values []string

	a, ok := e.Attributes[attr]
	if !ok {
		a = ldap.NewEntryAttribute(attr, value)
		values = value
	} else {
		oldValues := mapset.NewSetFromSlice(SliceToInterface(a.Values))
		newValues := mapset.NewSetFromSlice(SliceToInterface(value))

		valueSet := newValues.Difference(oldValues)

		for _, v := range valueSet.ToSlice() {
			a.Values = append(a.Values, v.(string))
			values = append(values, v.(string))
		}
	}

	if len(values) > 0 {
		e.Attributes[attr] = a
		e.AddAttributeChange("add", attr, values)
	}
}

func (e *Entry) AddAttributeValue(attr string, value string) {
	e.AddAttributeValues(attr, []string{value})
}

func (e *Entry) ReplaceAttributeValues(attr string, value []string) {
	if reflect.DeepEqual(value, e.Attributes[attr].Values) {
		return
	}
	e.Attributes[attr] = ldap.NewEntryAttribute(attr, value)
	e.AddAttributeChange("replace", attr, value)
}

func (e *Entry) ReplaceAttributeValue(attr string, value string) {
	e.ReplaceAttributeValues(attr, []string{value})
}

func (e *Entry) DeleteAttributeValues(attr string, value []string) {
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

func (e *Entry) DeleteAttributeValue(attr string, value string) {
	e.DeleteAttributeValues(attr, []string{value})
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
