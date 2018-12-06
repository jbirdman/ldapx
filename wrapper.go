package ldapx

import "gopkg.in/ldap.v2"

func NewModifyRequest(dn string) *ldap.ModifyRequest {
	return ldap.NewModifyRequest(dn)
}

func NewPasswordModifyRequesy(userIdentity string, oldPassword string, newPassword string) *ldap.PasswordModifyRequest {
	return ldap.NewPasswordModifyRequest(userIdentity, oldPassword, newPassword)
}

func NewSearchRequest(
	BaseDN string,
	Scope, DerefAliases, SizeLimit, TimeLimit int,
	TypesOnly bool,
	Filter string,
	Attributes []string,
	Controls []ldap.Control,
) *ldap.SearchRequest {
	return ldap.NewSearchRequest(BaseDN, Scope, DerefAliases, SizeLimit, TimeLimit, TypesOnly, Filter, Attributes, Controls)
}
