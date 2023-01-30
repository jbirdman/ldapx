package ldapx

import (
	"github.com/go-ldap/ldap/v3"
)

// LDAPSchema represents the LDAP schema.
type LDAPSchema struct {
	Syntaxes        []string // Attribute syntaxes
	MatchingRules   []string // Attribute matching rules
	MatchingRuleUse []string // Attribute matching rule use
	AttributeTypes  []string // Attribute types
	ObjectClasses   []string // Object classes
}

// RootDSE represents the RootDSE.
type RootDSE struct {
	SupportedLDAPVersion         []string // LDAP versions supported by the server
	SupportedControls            []string // Controls supported by the server
	SupportedExtensions          []string // Extensions supported by the server
	SupportedFeatures            []string // Features supported by the server
	SupportedSASLMechanisms      []string // SASL mechanisms supported by the server
	SupportedTLSCiphers          []string // TLS ciphers supported by the server
	ConfigContext                string   // Config context
	NamingContexts               []string // Naming contexts
	SubschemaSubEntry            string   // Subschema subentry
	SupportedAuthPasswordSchemes []string // Password schemes supported by the server
	VendorName                   string   // Vendor name
	VendorVersion                string   // Vendor version
}

// AttributeType represents an attribute type.
type AttributeType struct {
	OID         string
	Name        string
	Syntax      string
	SingleValue bool
}

// Schema returns the LDAP schema.
func (c *Conn) Schema() (*LDAPSchema, error) {
	rootDSE, err := c.RootDSE()
	if err != nil {
		return nil, err
	}

	//
	result, err := c.Search(NewSearchRequest(
		rootDSE.SubschemaSubEntry,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases,
		0, 0, false,
		"(objectclass=*)",
		[]string{
			"ldapSyntaxes",
			"matchingRules",
			"matchingRuleUse",
			"attributeTypes",
			"objectClasses",
		},
		nil,
	))
	if err != nil {
		return nil, err
	}

	for _, e := range result.Entries {
		return &LDAPSchema{
			Syntaxes:        e.GetAttributeValues("ldapSyntaxes"),
			MatchingRules:   e.GetAttributeValues("matchingRules"),
			MatchingRuleUse: e.GetAttributeValues("matchingRuleUse"),
			AttributeTypes:  e.GetAttributeValues("attributeTypes"),
			ObjectClasses:   e.GetAttributeValues("objectClasses"),
		}, nil
	}
	return nil, err
}

// rootDSE returns the RootDSE.
func rootDSE(conn *ldap.Conn) (*RootDSE, error) {
	result, err := conn.Search(NewSearchRequest(
		"",
		ldap.ScopeBaseObject, ldap.NeverDerefAliases,
		0, 0, false,
		"(objectclass=*)",
		[]string{
			"supportedLDAPVersion",
			"supportedControl",
			"supportedExtension",
			"supportedFeatures",
			"supportedSASLMechanisms",
			"supportedTLSCiphers",
			"configContext",
			"namingContexts",
			"subschemaSubentry",
			"supportedAuthPasswordSchemes",
			"vendorName",
			"vendorVersion",
		},
		nil,
	))
	if err != nil {
		return nil, err
	}
	for _, e := range result.Entries {
		return &RootDSE{
			SupportedLDAPVersion:         e.GetAttributeValues("supportedLDAPVersion"),
			SupportedControls:            e.GetAttributeValues("supportedControl"),
			SupportedExtensions:          e.GetAttributeValues("supportedExtension"),
			SupportedFeatures:            e.GetAttributeValues("supportedFeatures"),
			SupportedSASLMechanisms:      e.GetAttributeValues("supportedSASLMechanisms"),
			SupportedTLSCiphers:          e.GetAttributeValues("supportedTLSCiphers"),
			ConfigContext:                e.GetAttributeValue("configContext"),
			NamingContexts:               e.GetAttributeValues("namingContexts"),
			SubschemaSubEntry:            e.GetAttributeValue("subschemaSubentry"),
			SupportedAuthPasswordSchemes: e.GetAttributeValues("supportedAuthPasswordSchemes"),
			VendorName:                   e.GetAttributeValue("vendorName"),
			VendorVersion:                e.GetAttributeValue("vendorVersion"),
		}, nil
	}
	return nil, nil
}

// RootDSE returns the RootDSE.
func (c *Conn) RootDSE() (*RootDSE, error) {
	conn, err := getConn(c.pool)
	if err != nil {
		return nil, err
	}
	defer putConn(c.pool, conn)
	return rootDSE(conn)
}
