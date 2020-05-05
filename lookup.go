package ldapx

import "github.com/go-ldap/ldap/v3"

func (c Conn) Lookup(dn string) (*Entry, error) {
	result, err := c.Search(NewSearchRequest(dn, ldap.ScopeBaseObject, ldap.DerefAlways, 1, 0, false, "(objectclass=*)", nil, nil))
	if err != nil {
		return nil, err
	}

	return NewEntryFromLdapEntry(result.Entries[0]),nil
}
