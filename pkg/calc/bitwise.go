package calc

func IsBitSet(b byte, bit int) bool {
	return b >> bit & 1 == 1
}

