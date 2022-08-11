//go:build integration
// +build integration

package ldapx

import (
	"crypto/tls"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
)

const (
	server       = "ldaps://localhost:1636"
	binddn       = ""
	bindpassword = ""
)

var tlsConfig = &tls.Config{InsecureSkipVerify: true}

func TestOpenURL(t *testing.T) {
	_, err := OpenURL(server, binddn, bindpassword, tlsConfig)

	assert.Nil(t, err)
}

func TestConn_Search(t *testing.T) {
	conn, err := OpenURL(server, binddn, bindpassword, tlsConfig)

	assert.NoError(t, err, "Connection error")

	request := NewSearchRequest("dc=jcu,dc=edu,dc=au", ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(dc=jcu)", nil, nil)
	result, err := conn.Search(request)

	assert.NoError(t, err, "Search")
	assert.EqualValues(t, 1, len(result.Entries))
}

func TestConn_SearchWithPaging(t *testing.T) {
	conn, err := OpenURL(server, binddn, bindpassword, tlsConfig)

	assert.NoError(t, err, "Connection error")

	request := NewSearchRequest("ou=users,dc=jcu,dc=edu,dc=au", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(uid=test*)", nil, nil)
	result, err := conn.SearchWithPaging(request, 100)

	assert.NoError(t, err, "Search")
	assert.Condition(t, func() bool { return len(result.Entries) > 0 })
}

func TestConn_CheckBind(t *testing.T) {
	conn, err := OpenURL(server, binddn, bindpassword, tlsConfig)

	assert.NoError(t, err, "Connection error")

	err = conn.CheckBind("uid=testuser,ou=users,dc=jcu,dc=edu,dc=au", "testpassword")
	assert.NoError(t, err)
}

func TestConn_Compare(t *testing.T) {
	conn, err := OpenURL(server, binddn, bindpassword, tlsConfig)

	assert.NoError(t, err, "Connection error")

	result, err := conn.Compare("uid=testuser,ou=users,dc=jcu,dc=edu,dc=au", "uid", "testuser")
	assert.True(t, result, "Check for matching value")

	result, err = conn.Compare("uid=testuser,ou=users,dc=jcu,dc=edu,dc=au", "uid", "testuser1")
	assert.False(t, result, "Check for non-matching value")
}
