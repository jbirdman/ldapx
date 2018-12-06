package ldapx

func SliceToInterface(t []string) []interface{} {
	s := make([]interface{}, len(t))

	for i, v := range t {
		s[i] = v
	}

	return s
}

func InterfaceToSlice(t []interface{}) []string {
	s := make([]string, len(t))

	for i, v := range t {
		s[i] = v.(string)
	}

	return s
}

