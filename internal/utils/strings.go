package utils

// FixFormat removes the first and the last character of string.
// This solution works even for non-unicode.
func FixFormat(s string) string {
	return trimFLastRune(trimFirstRune(s))
}

func trimFirstRune(s string) string {
	for i := range s {
		if i > 0 {
			// The value i is the index in s of the second
			// rune.  Slice to remove the first rune.
			return s[i:]
		}
	}
	// There are 0 or 1 runes in the string.
	return ""
}

func trimFLastRune(s string) string {
	for i := range s {
		if i > 0 {
			return s[:1]
		}
	}
	// There are 0 or 1 runes in the string.
	return ""

}