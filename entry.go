package ldapx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"sort"
	"strings"
)

const (
	ChangeAdd    = "add"
	ChangeUpdate = "update"
	ChangeDelete = "delete"
)

type Entry struct {
	DN                 string            `json:"dn"`
	ChangeType         string            `json:"change_type"`
	Attributes         AttributeMap      `json:"attributes,omitifempty"`
	Changes            []AttributeChange `json:"changes,omitifempty"`
	committed          bool
	originalAttributes AttributeMap
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
	return &Entry{
		DN:                 dn,
		Attributes:         NewAttributeMap(),
		originalAttributes: NewAttributeMap(),
		ChangeType:         ChangeAdd}
}

func NewEntryFromLdapEntry(entry *ldap.Entry) *Entry {
	e := NewEntry("")

	if entry != nil {
		e.DN = entry.DN
		e.ChangeType = ChangeUpdate

		// Copy in the attributes from the ldap entry
		for _, a := range entry.Attributes {
			e.Attributes.PutEntryAttribute(a)
			e.originalAttributes.PutEntryAttribute(a)
		}
	}

	return e
}

func (e *Entry) ToLdapEntry() *ldap.Entry {
	attrs := make(map[string][]string)

	for _, k := range e.AttributeNames() {
		attrs[k] = e.GetAttributeValues(k)
	}

	return ldap.NewEntry(e.DN, attrs)
}

func (e Entry) Print() {
	e.ToLdapEntry().Print()
	//fmt.Printf("DN: %s\n", e.DN)
	//for _, attr := range e.Attributes {
	//	attr.Print()
	//}
}

func (e Entry) PrettyPrint(indent int) {
	e.ToLdapEntry().PrettyPrint(indent)
}

func (e Entry) AttributeNames() []string {
	return e.Attributes.AttributeNames()
}

func (e Entry) ToLDIF() string {
	var b strings.Builder
	fmt.Fprintf(&b, "dn: %s\n", e.DN)

	keys := e.AttributeNames()
	sort.Strings(keys)

	for _, k := range keys {
		for _, v := range e.Attributes.Get(k).Values {
			fmt.Fprintf(&b, "%s: %s\n", k, v)
		}
	}
	return b.String()
}

func (e Entry) Changed() bool {
	return len(e.Changes) > 0
}

func (e *Entry) ResetChanges() {
	e.Changes = nil
}

func (e *Entry) GetAttributeValues(attribute string) []string {
	v := e.Attributes.Get(attribute)
	if v == nil {
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

func (e Entry) Clone() *Entry {
	dest := NewEntry(e.DN)

	for _, a := range e.AttributeNames() {
		dest.AddAttributeValues(a, e.GetAttributeValues(a))
	}

	return dest
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

func (e *Entry) RenameAttribute(from, to string) {
	e.Attributes.Rename(from, to)
}

func (e Entry) ToJSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}
