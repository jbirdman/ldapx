package entities

import (
	"fmt"

	"github.com/jbirdman/ldapx"
)

func NewOrganizationalUnit(baseDN, ou string) *ldapx.Entry {
	entry := ldapx.NewEntry(fmt.Sprintf("ou=%s,%s", ou, baseDN))

	entry.ReplaceAttributeValues("objectclass", []string{"top", "organizationalUnit"})
	entry.ReplaceAttributeValue("ou", ou)

	return entry
}
