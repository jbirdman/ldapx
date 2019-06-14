package ldapx

import "gopkg.in/ldap.v2"

func (e *Entry) DeleteAttributeValues(attr string, value []string) {
	old, ok := e.Attributes[attr]
	if !ok {
		return
	}
	if len(value) > 0 {
		a := ldap.NewEntryAttribute(attr, nil)

		for _, v := range old.Values {
			found := false
			for _, nv := range value {
				if nv == v {
					found = true
					break
				}
			}

			if !found {
				a.Values = append(a.Values, v)
			}
		}
		e.Attributes[attr] = a
	} else {
		delete(e.Attributes, attr)
	}

	e.AddAttributeChange("delete", attr, value)
}

func (e *Entry) DeleteAttributeValue(attr string, value string) {
	e.DeleteAttributeValues(attr, []string{value})
}

