package format

// Sections is a super-class of a basic NES ROM section. To extract those section areas you need to use a rom.Parser().
// There are different parsing implementations, found at rom/format/ directory (e.g. iNES, or iNES 2.0).
type Sections struct {
	Header            []byte	// Added by a person, either iNES or iNES 2.0. Required by emulators.
	Trainer           []byte	// Hacks and stuff
	PRGROM            []byte	// Memory chip connected to the CPU. Contains the code.
	CHRROM            []uint8	// Memory chip connected to the PPU. Contains a fixed set of graphics tile data.
	PRGRAM            []byte	// Rare: There may be an additional chip like that to hold even more data.
	CHRRAM			  []byte	// Rare: Some cartridges have this chip to hold data the CPU has copied from PRGROM.
	PlayChoiceInstRom []byte
	PlayChoicePROM    []byte
}