package ldapx

import "github.com/go-ldap/ldap/v3"

type LDAPSchema struct {
	Syntaxes        []string
	MatchingRules   []string
	MatchingRuleUse []string
	AttributeTypes  []string
	ObjectClasses   []string
}

type RootDSE struct {
	SupportedLDAPVersion         []string
	SupportedControls            []string
	SupportedExtensions          []string
	SupportedFeatures            []string
	SupportedSASLMechanisms      []string
	SupportedTLSCiphers          []string
	ConfigContext                string
	NamingContexts               []string
	SubschemaSubEntry            string
	SupportedAuthPasswordSchemes []string
	VendorName                   string
	VendorVersion                string
}

func (c *Conn) Schema() (*LDAPSchema, error) {
	rootDSE, err := c.RootDSE()
	if err != nil {
		return nil, err
	}

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

func (c *Conn) RootDSE() (*RootDSE, error) {
	result, err := c.Search(NewSearchRequest(
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
