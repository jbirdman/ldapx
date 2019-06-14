package ldapx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntry(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.AddAttributeValues("cn", []string{"test"})
	entry.AddAttributeValues("objectclass", []string{"top", "person"})

	assert.Equal(t, "cn=test",entry.DN)
	assert.Equal(t, "test", entry.GetAttributeValue("cn"))
}

func TestEntry_AddAttribute(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttributeValues("cn", []string{"test"})

	assert.Equal(t, "test", entry.GetAttributeValue("cn"))
}

func TestEntry_ReplaceAttribute(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttributeValues("cn", []string{"test"})
	entry.ReplaceAttributeValues("cn", []string{"test2"})

	assert.Equal(t, "test2", entry.GetAttributeValue("cn"))
}

func TestEntry_ReplaceAttributeIdempotent(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttributeValues("cn", []string{"test"})
	entry.ResetChanges()
	entry.ReplaceAttributeValues("cn", []string{"test"})

	assert.Equal(t, "test", entry.GetAttributeValue("cn"))
	assert.Equal(t,0,len(entry.Changes))
}

func TestEntry_AddAttributeValuesIdempotent(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttributeValues("cn", []string{"test"})
	entry.ResetChanges()
	entry.AddAttributeValues("cn", []string{"test","test1"})

	assert.Equal(t, []string{"test","test1"}, entry.GetAttributeValues("cn"))
	assert.Equal(t,1,len(entry.Changes))
}
