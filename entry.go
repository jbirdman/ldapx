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

func NewEntry(entry *ldap.Entry) Entry {
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

	return Entry{ChangeType: changeType, DN: dn, Attributes: attributes}
}

func (c *Entry) Print() {
	fmt.Printf("DN: %s\n", c.DN)
	for _, attr := range c.Attributes {
		attr.Print()
	}
}

func (c *Entry) Clear() {
	c.ChangeType = CHANGE_ADD
	c.Changes = nil
}

func (c *Entry) GetAttributeValues(attribute string) []string {
	return c.Attributes[attribute].Values
}

func (e *Entry) GetAttributeValue(attribute string) string {
	values := e.GetAttributeValues(attribute)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *Entry) AddAttributeChange(action string, attr string, value []string) {
	c.Changes = append(c.Changes, AttributeChange{Action: action, Attr: attr, Value: value})
}
func (c *Entry) AddAttribute(attr string, value []string) {
	a, ok := c.Attributes[attr]
	if !ok {
		a = ldap.NewEntryAttribute(attr, value)
	} else {
		a.Values = append(a.Values, value...)
	}
	c.Attributes[attr] = a
	c.AddAttributeChange("add", attr, value)
}

func (c *Entry) ReplaceAttribute(attr string, value []string) {
	c.Attributes[attr] = ldap.NewEntryAttribute(attr, value)
	c.AddAttributeChange("replace", attr, value)
}

func (c *Entry) DeleteAttribute(attr string, value []string) {
	old, ok := c.Attributes[attr]
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
		c.Attributes[attr] = a
	} else {
		delete(c.Attributes, attr)
	}

	c.AddAttributeChange("delete", attr, value)
}

func (c Entry) Update(conn *Conn) error {
	if len(c.Changes) == 0 {
		return nil
	}

	switch c.ChangeType {
	case CHANGE_ADD:
		return conn.Add(buildAddRequest(c.DN, c.Changes))
	case CHANGE_UPDATE:
		return conn.Modify(buildModifyRequest(c.DN, c.Changes))
	case CHANGE_DELETE:
		return conn.Del(buildDelRequest(c.DN))
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
