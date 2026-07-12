package cmd

import "strings"

func isValidShellIdentifier(name string) bool {
	if name == "" {
		return false
	}

	for i, r := range name {
		if i == 0 {
			if !isShellIdentifierStart(r) {
				return false
			}
			continue
		}

		if !isShellIdentifierPart(r) {
			return false
		}
	}

	return true
}

func isShellIdentifierStart(r rune) bool {
	return r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isShellIdentifierPart(r rune) bool {
	return isShellIdentifierStart(r) || (r >= '0' && r <= '9')
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
