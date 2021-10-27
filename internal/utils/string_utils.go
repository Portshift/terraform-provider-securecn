package utils

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FilterEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func IsStringSlicesIdentical(s1 []string, s2 []string) bool {
	if s1 == nil && s2 != nil {
		return false
	}

	if s2 == nil && s1 != nil {
		return false
	}

	if len(s1) != len(s2) {
		return false
	}

	for _, str := range s1 {
		if !ContainsString(s2, str) {
			return false
		}
	}

	for _, str := range s2 {
		if !ContainsString(s1, str) {
			return false
		}
	}
	return true
}
