package rom

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/calc"
	"github.com/drpaneas/sengo/pkg/rom/cartridge"
	"github.com/drpaneas/sengo/pkg/rom/format"
	"github.com/drpaneas/sengo/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type Rom struct {
	Name   string
	File   Properties
	Parser parser
	Data   cartridge.ContentInSection
}

func (r *Rom) Parse(rom []byte) cartridge.ContentInSection {
	return r.Parser.Parse(rom)
}

func Open(romFilepath string) Rom {
	// read the file
	bytes := utils.ReadRom(romFilepath)

	// Initialize data structures
	data := ContentInType{
		Binary: calc.BinaryDump(bytes),
		Hex:    calc.Hexdump(bytes),
		Bytes:  bytes,
		ASCII:  calc.AsciiDump(calc.Hexdump(bytes)),
	}

	file := Properties{
		Path: romFilepath,
		Name: filepath.Base(romFilepath),
		Size: len(bytes),
		Data: data,
	}

	rom := Rom{
		Name:   strings.TrimSuffix(file.Name, filepath.Ext(file.Name)),	// remove the *.nes suffix
		File:   file,
		Parser: selectParser(data), // strategy pattern
		Data:   cartridge.ContentInSection{},
	}

	if rom.Parser == nil {
		fmt.Println("Error: Couldn't parse the ROM (unknown format)")
		os.Exit(1)
	}

	return rom
}

func isINES2(header ContentInType) bool {
	return calc.IsBitSet(header.Bytes[7], 3) && !calc.IsBitSet(header.Bytes[7], 2)
}

func isINES(header ContentInType) bool {
	if header.ASCII[0] == "N" && header.ASCII[1] == "E" && header.ASCII[2] == "S" {
		return header.Hex[3] == "1A" // this is the <EOF> equivalent in Hex
	}
	return false
}

func selectParser(header ContentInType) parser {
	if isINES(header) {
		if isINES2(header) {
			return format.INES2{}
		} else {
			return format.INES{}
		}
	}
	return nil
}
