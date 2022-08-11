package ldapx

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type DN ldap.DN

func ParseDN(str string) (*DN, error) {
	dn, err := ldap.ParseDN(str)
	return (*DN)(dn), err
}

func (dn *DN) ToString() string {
	var dnc = make([]string, 0, len(dn.RDNs))
	for _, c := range dn.RDNs {
		dnc = append(dnc, joinRDNAttrs(c))

	}
	return strings.Join(dnc, ",")
}

func joinRDNAttrs(rdn *ldap.RelativeDN) string {
	var attrs = make([]string, 0, len(rdn.Attributes))
	for _, a := range rdn.Attributes {
		attrs = append(attrs, fmt.Sprintf("%s=%s", a.Type, a.Value))
	}
	return strings.Join(attrs, "+")
}

func (dn *DN) Append(attr, value string) {
	dn.RDNs = append([]*ldap.RelativeDN{&ldap.RelativeDN{
		Attributes: []*ldap.AttributeTypeAndValue{&ldap.AttributeTypeAndValue{
			Type:  attr,
			Value: value,
		}},
	}}, dn.RDNs...)
}
