package ldapx

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

const (
	ChangeAdd    = "add"    // ChangeAdd is the change type for adding an entry
	ChangeUpdate = "update" // ChangeUpdate is the change type for updating an entry
	ChangeDelete = "delete" // ChangeDelete is the change type for deleting an entry
)

// Entry represents an LDAP entry
type Entry struct {
	DN                 string            `json:"dn"`                   // DN is the distinguished name of the entry
	ChangeType         string            `json:"change_type"`          // ChangeType is the type of change to be applied to the entry
	Attributes         AttributeMap      `json:"attributes,omitempty"` // Attributes is a map of attribute name to attribute
	Changes            []AttributeChange `json:"changes,omitempty"`    // Changes is a list of changes to be applied to the entry
	committed          bool              // committed is true if the entry has been committed to the server
	originalAttributes AttributeMap      // originalAttributes is a copy of the attributes when the entry was created
}

var _ MutableEntry = &Entry{}

// MutableEntry is an interface for an entry that can be modified
type MutableEntry interface {
	AddAttributeValueIgnoreCase(attr string, value string, ignoreCase bool)        // AddAttributeValueIgnoreCase adds a value to an attribute, ignoring case if ignoreCase is true
	AddAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool)     // AddAttributeValuesIgnoreCase adds values to an attribute, ignoring case if ignoreCase is true
	AddAttributeValue(attr string, value string)                                   // AddAttributeValue adds a value to an attribute
	AddAttributeValues(attr string, value []string)                                // AddAttributeValues adds values to an attribute
	ReplaceAttributeValueIgnoreCase(attr string, value string, ignoreCase bool)    // ReplaceAttributeValueIgnoreCase replaces an attribute value, ignoring case if ignoreCase is true
	ReplaceAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool) // ReplaceAttributeValuesIgnoreCase replaces attribute values, ignoring case if ignoreCase is true
	ReplaceAttributeValue(attr string, value string)                               // ReplaceAttributeValue replaces an attribute value
	ReplaceAttributeValues(attr string, value []string)                            // ReplaceAttributeValues replaces attribute values
	DeleteAttributeValueIgnoreCase(attr string, value string, ignoreCase bool)     // DeleteAttributeValueIgnoreCase deletes an attribute value, ignoring case if ignoreCase is true
	DeleteAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool)  // DeleteAttributeValuesIgnoreCase deletes attribute values, ignoring case if ignoreCase is true
	DeleteAttributeValue(attr string, value string)                                // DeleteAttributeValue deletes an attribute value
	DeleteAttributeValues(attr string, value []string)                             // DeleteAttributeValues deletes attribute values
	DeleteAttribute(attr string)                                                   // DeleteAttribute deletes an attribute
	SyncAttributeValues(attr string, value []string)                               // SyncAttributeValues syncs an attribute with the given values
	SyncAttributeValuesIgnoreCase(attr string, value []string, ignoreCase bool)    // SyncAttributeValuesIgnoreCase syncs an attribute with the given values, ignoring case if ignoreCase is true
	Update(conn *Conn) error                                                       // Update updates the entry on the server
	Changed() bool                                                                 // Changed returns true if the entry has been changed
}

// AttributeChange represents a change to an attribute
type AttributeChange struct {
	Action string
	Attr   string
	Value  []string
}

// convert string slice to an interface slice

// NewEntry creates a new entry with the given DN
func NewEntry(dn string) *Entry {
	return &Entry{
		DN:                 dn,
		Attributes:         NewAttributeMap(),
		originalAttributes: NewAttributeMap(),
		ChangeType:         ChangeAdd,
	}
}

// NewEntryFromLdapEntry creates a new entry from a ldap entry
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

// ToLdapEntry converts the ldapx entry to a ldap entry
func (e *Entry) ToLdapEntry() *ldap.Entry {
	// Copy the attributes to a map
	attrs := make(map[string][]string)

	// Copy in the attributes from the ldap entry
	for _, k := range e.AttributeNames() {
		attrs[k] = e.GetAttributeValues(k)
	}

	// Create the ldap entry
	return ldap.NewEntry(e.DN, attrs)
}

func (e *Entry) Print() {
	// Print the entry
	e.ToLdapEntry().Print()
}

// PrettyPrint prints the entry with the given indent
func (e *Entry) PrettyPrint(indent int) {
	// Print the entry
	e.ToLdapEntry().PrettyPrint(indent)
}

// AttributeNames returns the names of the attributes
func (e *Entry) AttributeNames() []string {
	// Return the attribute names
	return e.Attributes.AttributeNames()
}

// ToLDIF converts the entry to LDIF
func (e *Entry) ToLDIF() string {
	// Create a string builder
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "dn: %s\n", e.DN)

	// Sort the attributes
	keys := e.AttributeNames()
	sort.Strings(keys)

	// Print the attributes
	for _, k := range keys {
		for _, v := range e.Attributes.Get(k).Values {
			_, _ = fmt.Fprintf(&b, "%s: %s\n", k, v)
		}
	}

	// Return the string
	return b.String()
}

func (e *Entry) Changed() bool {
	// Check if the entry has changed
	return len(e.Changes) > 0
}

func (e *Entry) ResetChanges() {
	// Reset the changes
	e.Changes = nil
}

func (e *Entry) AttributeExists(attr string) bool {
	// Check if the attribute exists
	return e.Attributes.AttributeExists(attr)
}

func (e *Entry) GetAttributeValues(attribute string) []string {
	// Get the attribute values
	v := e.Attributes.Get(attribute)
	if v == nil {
		return nil
	}
	return v.Values
}

func (e *Entry) GetAttributeValue(attribute string) string {
	// Get the attribute value
	values := e.GetAttributeValues(attribute)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (e *Entry) AddAttributeChange(action string, attr string, value []string) {
	// Add the attribute change
	e.Changes = append(e.Changes, AttributeChange{Action: action, Attr: attr, Value: value})
}

func (e *Entry) Update(conn *Conn) error {
	// Update the entry
	if !e.Changed() {
		return nil
	}

	// Check if the entry has been committed
	if e.committed {
		return errors.New("entry can only be updated once")
	}
	e.committed = true

	switch e.ChangeType {
	case ChangeAdd:
		// Add the entry
		return conn.Add(buildAddRequest(e.DN, e.Changes))
	case ChangeUpdate:
		// Modify the entry
		return conn.Modify(buildModifyRequest(e.DN, e.Changes))
	case ChangeDelete:
		// Delete the entry
		return conn.Del(buildDelRequest(e.DN))
	}

	return nil
}

func (e *Entry) Clone() *Entry {
	// Clone the entry
	dest := NewEntry(e.DN)

	for _, a := range e.AttributeNames() {
		dest.AddAttributeValues(a, e.GetAttributeValues(a))
	}

	return dest
}

// buildAddRequest builds an add request from the changes
func buildAddRequest(dn string, changes []AttributeChange) *ldap.AddRequest {
	// Build the add request
	r := NewAddRequest(dn, nil)

	// Add the attributes
	for _, change := range changes {
		// Check if the attribute is being deleted
		if change.Action != "add" && change.Action != "replace" {
			// Skip the attribute
			continue
		}
		// Add the attribute
		r.Attribute(change.Attr, change.Value)
	}

	return r
}

// buildModifyRequest builds a modify request from the changes
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

// buildDelRequest builds a delete request
func buildDelRequest(dn string) *ldap.DelRequest {
	return NewDelRequest(dn, nil)
}

// RenameAttribute renames an attribute
func (e *Entry) RenameAttribute(from, to string) {
	e.Attributes.Rename(from, to)
}

// ToJSON returns the entry as a JSON string
func (e *Entry) ToJSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}
