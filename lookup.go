package ldapx

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func (c *Conn) Lookup(dn string) (*Entry, error) {
	result, err := c.Search(NewSearchRequest(dn, ldap.ScopeBaseObject, ldap.DerefAlways, 1, 0, false, "(objectclass=*)", nil, nil))
	if err != nil {
		return nil, err
	}

	return NewEntryFromLdapEntry(result.Entries[0]), nil
}

func (c *Conn) FindEntry(dn string, filter string, attributes []string) (*Entry, error) {
	return FindEntry(c, dn, filter, attributes)
}

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

func (c *Conn) QuickSearch(dn string, filter string, attributes []string) (*ldap.SearchResult, error) {
	return c.Search(NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false, filter, attributes, nil))
}
