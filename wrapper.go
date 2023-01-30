package ldapx

import "github.com/go-ldap/ldap/v3"

// NewModifyRequest creates a new ModifyRequest with the given DN and controls.
func NewModifyRequest(dn string, controls []ldap.Control) *ldap.ModifyRequest {
	return ldap.NewModifyRequest(dn, controls)
}

// NewPasswordModifyRequest creates a new PasswordModifyRequest with the given user identity, old password and new password.
func NewPasswordModifyRequest(userIdentity string, oldPassword string, newPassword string) *ldap.PasswordModifyRequest {
	return ldap.NewPasswordModifyRequest(userIdentity, oldPassword, newPassword)
}

// NewSearchRequest creates a new SearchRequest with the given base DN, scope,
// derefAliases, sizeLimit, timeLimit, typesOnly, filter, attributes and
// controls.
func NewSearchRequest(
	baseDN string,
	scope, derefAliases, sizeLimit, timeLimit int,
	typesOnly bool,
	filter string,
	attributes []string,
	controls []ldap.Control,
) *ldap.SearchRequest {
	return ldap.NewSearchRequest(baseDN, scope, derefAliases, sizeLimit, timeLimit, typesOnly, filter, attributes, controls)
}

// NewAddRequest creates a new AddRequest with the given DN and controls.
func NewAddRequest(dn string, controls []ldap.Control) *ldap.AddRequest {
	return ldap.NewAddRequest(dn, controls)
}

// NewDelRequest creates a new DelRequest with the given DN and controls.
func NewDelRequest(dn string, controls []ldap.Control) *ldap.DelRequest {
	return ldap.NewDelRequest(dn, controls)
}

// NewControlBeheraPasswordPolicy creates a new Behera password policy control.
func NewControlBeheraPasswordPolicy() *ldap.ControlBeheraPasswordPolicy {
	return ldap.NewControlBeheraPasswordPolicy()
}

// NewControlPaging creates a new paging control with the given paging size.
func NewControlPaging(pagingSize uint32) *ldap.ControlPaging {
	return ldap.NewControlPaging(pagingSize)
}

// NewControlString creates a new ControlString with the given control type, criticality and control value.
func NewControlString(controlType string, criticality bool, controlValue string) *ldap.ControlString {
	return ldap.NewControlString(controlType, criticality, controlValue)
}
