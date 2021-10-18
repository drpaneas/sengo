package cartridge

import "github.com/drpaneas/sengo/pkg/rom/section"

// ContentInSection is a super-class of a basic NES ROM section. To extract those section areas you need to use a rom.Parser().
// There are different parsing implementations, found at rom/format/ directory (e.g. iNES, or iNES 2.0).
type ContentInSection struct {
	Header            section.Header
	Trainer           section.Trainer
	PRGROM            section.ProgramRom
	CHRROM            section.CharacterRom
	PlayChoiceInstRom section.PlayChoiceInstRom
	PlayChoicePROM    section.PlayChoicePROM
}