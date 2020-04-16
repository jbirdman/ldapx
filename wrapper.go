package ldapx

import "github.com/go-ldap/ldap/v3"

func NewModifyRequest(dn string, controls []ldap.Control) *ldap.ModifyRequest {
	return ldap.NewModifyRequest(dn, controls)
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

func NewAddRequest(dn string, controls []ldap.Control) *ldap.AddRequest {
	return ldap.NewAddRequest(dn, controls)
}

func NewDelRequest(dn string, controls []ldap.Control) *ldap.DelRequest {
	return ldap.NewDelRequest(dn, controls)
}

func NewControlBeheraPasswordPolicy() *ldap.ControlBeheraPasswordPolicy {
	return ldap.NewControlBeheraPasswordPolicy()
}

func NewControlPaging(pagingSize uint32) *ldap.ControlPaging {
	return ldap.NewControlPaging(pagingSize)
}

func NewControlString(controlType string, criticality bool, controlValue string) *ldap.ControlString {
	return ldap.NewControlString(controlType, criticality, controlValue)
}
