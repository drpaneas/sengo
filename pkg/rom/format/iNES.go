package format

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/rom/cartridge"
	"github.com/drpaneas/sengo/pkg/rom/section"
)

// INES is a sub-class of section.ContentInSection
type INES struct {
	cartridge.ContentInSection
}

func (i INES) Parse(rom []byte) cartridge.ContentInSection {

	fmt.Println("This is iNES")

	return cartridge.ContentInSection{
		Header:            section.Header{},
		Trainer:           section.Trainer{},
		PRGROM:            section.ProgramRom{},
		CHRROM:            section.CharacterRom{},
		PlayChoiceInstRom: section.PlayChoiceInstRom{},
		PlayChoicePROM:    section.PlayChoicePROM{},
	}
}