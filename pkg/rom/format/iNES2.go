package format

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/rom/cartridge"
	"github.com/drpaneas/sengo/pkg/rom/section"
)

// INES2 is a sub-class of section.ContentInSection
type INES2 struct {
	cartridge.ContentInSection
}

func (i INES2) Parse(rom []byte) cartridge.ContentInSection {

	fmt.Println("This is INES 2.0")

	return cartridge.ContentInSection{
		Header:            section.Header{},
		Trainer:           section.Trainer{},
		PRGROM:            section.ProgramRom{},
		CHRROM:            section.CharacterRom{},
		PlayChoiceInstRom: section.PlayChoiceInstRom{},
		PlayChoicePROM:    section.PlayChoicePROM{},
	}
}
