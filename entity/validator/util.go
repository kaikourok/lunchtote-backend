package validator

import "unicode"

func isOnlySpaces(s string) bool {
	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}

	return true
}

func isContainSpecialRunes(s string) bool {
	for _, c := range s {
		if unicode.In(c, unicode.C) && c != '\r' && c != '\n' {
			return true
		}
	}

	return false
}

func isHexadecimal(s string) bool {
	for _, r := range s {
		if !('0' <= r && r <= '9') && !('A' <= r && r <= 'F') && !('a' <= r && r <= 'f') {
			return false
		}
	}

	return true
}
