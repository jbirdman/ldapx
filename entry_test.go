package ldapx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntry(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.AddAttribute("cn", []string{"test"})
	entry.AddAttribute("objectclass", []string{"top", "person"})

	entry.Print()

	r := buildAddRequest(entry.DN, entry.Changes)
	fmt.Printf("request: %+v\n", r)

	assert.True(t, true)
}

func TestEntry_AddAttribute(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttribute("cn", []string{"test"})

	assert.Equal(t, "test", entry.GetAttributeValue("cn"))
}

func TestEntry_ReplaceAttribute(t *testing.T) {
	entry := NewEntry("cn=test")

	entry.DN = "cn=test"
	entry.AddAttribute("cn", []string{"test"})
	entry.ReplaceAttribute("cn", []string{"test2"})

	fmt.Println(entry.Changes)

	assert.Equal(t, "test2", entry.GetAttributeValue("cn"))
}
