package ldapx

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// DN is a distinguished name.
type DN ldap.DN

// ParseDN parses the given string into a DN.
func ParseDN(str string) (*DN, error) {
	dn, err := ldap.ParseDN(str)
	return (*DN)(dn), err
}

// ToString returns the string representation of the DN.
func (dn *DN) ToString() string {
	dnc := make([]string, 0, len(dn.RDNs))
	for _, c := range dn.RDNs {
		dnc = append(dnc, joinRDNAttrs(c))
	}
	return strings.Join(dnc, ",")
}

// joinRDNAttrs joins the attributes of the RDN into a string.
func joinRDNAttrs(rdn *ldap.RelativeDN) string {
	attrs := make([]string, 0, len(rdn.Attributes))
	for _, a := range rdn.Attributes {
		attrs = append(attrs, fmt.Sprintf("%s=%s", a.Type, a.Value))
	}
	return strings.Join(attrs, "+")
}

// Append appends the given attribute and value to the DN.
func (dn *DN) Append(attr, value string) {
	dn.RDNs = append([]*ldap.RelativeDN{{
		Attributes: []*ldap.AttributeTypeAndValue{{
			Type:  attr,
			Value: value,
		}},
	}}, dn.RDNs...)
}

// IsDN returns true if the string is a DN.
func IsDN(dn string) bool {
	_, err := ParseDN(dn)

	return err == nil
}
