package calc

func IsBitSet(b byte, bit int) bool {
	return b >> bit & 1 == 1
}

func ReadHighNibbleByte(B byte) byte {
	B = B >> 4
	if B > 15 {
		B = 15
	}
	return B
}

func ReadLowNibbleByte(B byte) byte {
	B = B << 4
	B = B >> 4
	if B > 15 {
		B = 15
	}
	return B
}

func MergeNibbles(highNibble byte, lowNibble byte) byte {
	highNibble = highNibble << 4
	return highNibble | lowNibble
}