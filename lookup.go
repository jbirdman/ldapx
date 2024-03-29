package ldapx

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// Lookup searches for the given DN and returns the entry.
func (c *Conn) Lookup(dn string) (*Entry, error) {
	result, err := c.Search(NewSearchRequest(dn, ldap.ScopeBaseObject, ldap.DerefAlways, 1, 0, false, "(objectclass=*)", nil, nil))
	if err != nil {
		return nil, err
	}

	return NewEntryFromLdapEntry(result.Entries[0]), nil
}

// LookupOrNew searches for the given DN and returns the entry if found, otherwise a new entry is created with the given DN.
func (c *Conn) LookupOrNew(dn string) (*Entry, error) {
	entry, err := c.Lookup(dn)
	if err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return nil, err
		}
		entry = NewEntry(dn)
	}
	return entry, nil
}

// FindEntry searches using given DN base and returns the first entry that matches the filter.
func (c *Conn) FindEntry(dn string, filter string, attributes []string) (*Entry, error) {
	return FindEntry(c, dn, filter, attributes)
}

// FindEntry searches using given DN base and returns the first entry that matches the filter using the given connection.
func FindEntry(conn *Conn, dn string, filter string, attributes []string) (*Entry, error) {
	result, err := conn.QuickSearch(dn, filter, attributes)
	if err != nil {
		return nil, err
	}

	// No entries found
	if len(result.Entries) == 0 {
		return nil, nil
	}

	// More than one entry matched
	if len(result.Entries) > 1 {
		return nil, fmt.Errorf("multiple entries matched")
	}

	return NewEntryFromLdapEntry(result.Entries[0]), nil
}

// QuickSearch performs a search using the given DN base, filter and attributes.
func (c *Conn) QuickSearch(dn string, filter string, attributes []string) (*ldap.SearchResult, error) {
	return c.Search(NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false, filter, attributes, nil))
}

// GetAttributeFromDN return value of first (leftmost) RDN that matches attribute name
func GetAttributeFromDN(attr, dn string) (string, error) {
	d, err := ldap.ParseDN(dn)
	if err != nil {
		return "", err
	}

	for _, a := range d.RDNs {
		for _, av := range a.Attributes {
			if strings.EqualFold(av.Type, attr) {
				return av.Value, nil
			}
		}
	}

	return "", fmt.Errorf("dn does not contain attribute '%s'", attr)
}
