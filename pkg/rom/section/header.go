package section

type Header struct {
	Data       []byte
	Signature  []byte
	Byte4 [1]byte
	Byte5 [1]byte
	Flags6 [1]byte
	Flags7 [1]byte
	Flags8 [1]byte
	Flags9 [1]byte
	Flags10 [1]byte
	UnusedPadding [5]byte
}