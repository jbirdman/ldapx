package entities

import (
	"fmt"
	"git.jcu.edu.au/go/ldapx"
)

func NewOrganizationalUnit(baseDN, ou string) *ldapx.Entry {
	entry := ldapx.NewEntry(fmt.Sprintf("ou=%s,%s", ou, baseDN))

	entry.ReplaceAttributeValues("objectclass", []string{"top", "organizationalUnit"})
	entry.ReplaceAttributeValue("ou", ou)

	return entry
}
