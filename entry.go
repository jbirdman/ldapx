package ldapx

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

const (
	ChangeAdd    = "add"
	ChangeUpdate = "update"
	ChangeDelete = "delete"
)

type Entry struct {
	ChangeType         string
	DN                 string
	Attributes         map[string]*ldap.EntryAttribute
	Changes            []AttributeChange
	committed          bool
	originalAttributes map[string]*ldap.EntryAttribute
}

type MutableEntry interface {
	AddAttributeValue(attr string, value string)
	AddAttributeValues(attr string, value []string)
	ReplaceAttributeValue(attr string, value string)
	ReplaceAttributeValues(attr string, value []string)
	DeleteAttributeValue(attr string, value string)
	DeleteAttributeValues(attr string, value []string)
	DeleteAttribute(attr string)
	SyncAttributeValues(attr string, value []string)
	Update(conn *Conn) error
	Changed() bool
}

type AttributeChange struct {
	Action string
	Attr   string
	Value  []string
}

func cloneAttribute(a *ldap.EntryAttribute) *ldap.EntryAttribute {
	return &ldap.EntryAttribute{
		Name:       a.Name,
		Values:     append([]string(nil), a.Values...),
		ByteValues: append([][]byte(nil), a.ByteValues...),
	}
}

func NewEntry(dn string) *Entry {
	return &Entry{DN: dn, Attributes: make(map[string]*ldap.EntryAttribute), originalAttributes: make(map[string]*ldap.EntryAttribute), ChangeType: ChangeAdd}
}

func NewEntryFromLdapEntry(entry *ldap.Entry) *Entry {
	e := NewEntry("")

	if entry != nil {
		e.DN = entry.DN
		e.ChangeType = ChangeUpdate

		// Copy in the attributes from the ldap entry
		for _, a := range entry.Attributes {
			e.Attributes[a.Name] = a
			e.originalAttributes[a.Name] = cloneAttribute(a)
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

func (e Entry) Print() {
	fmt.Printf("DN: %s\n", e.DN)
	for _, attr := range e.Attributes {
		attr.Print()
	}
}

func (e Entry) Changed() bool {
	return len(e.Changes) > 0
}

func (e *Entry) ResetChanges() {
	e.Changes = nil
}

func (e *Entry) GetAttributeValues(attribute string) []string {
	v, ok := e.Attributes[attribute]
	if !ok {
		return nil
	}
	return v.Values
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

func (e *Entry) Update(conn *Conn) error {
	if !e.Changed() {
		return nil
	}

	if e.committed {
		return errors.New("entry can only be updated once")
	}
	e.committed = true

	switch e.ChangeType {
	case ChangeAdd:
		return conn.Add(buildAddRequest(e.DN, e.Changes))
	case ChangeUpdate:
		return conn.Modify(buildModifyRequest(e.DN, e.Changes))
	case ChangeDelete:
		return conn.Del(buildDelRequest(e.DN))
	}

	return nil
}

func buildAddRequest(dn string, changes []AttributeChange) *ldap.AddRequest {
	r := NewAddRequest(dn, nil)

	for _, change := range changes {
		if change.Action != "add" && change.Action != "replace" {
			continue
		}

		r.Attribute(change.Attr, change.Value)
	}

	return r
}

func buildModifyRequest(dn string, changes []AttributeChange) *ldap.ModifyRequest {
	r := NewModifyRequest(dn, nil)

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
