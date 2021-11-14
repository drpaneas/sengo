package rom

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/calc"
	"github.com/drpaneas/sengo/pkg/rom/format"
	"github.com/drpaneas/sengo/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type Rom struct {
	Name    string
	File    Properties
	Parser  parser
	Content format.Sections
}

func newRom(data ContentType, file Properties) Rom {
	rom := Rom{
		Name:    strings.TrimSuffix(file.Name, filepath.Ext(file.Name)), // remove the *.nes suffix
		File:    file,
		Parser:  selectParser(data), // strategy pattern
		Content: format.Sections{},  // initialize
	}

	return rom
}

func Open(romFilepath string) Rom {
	// read the file
	bytes := utils.ReadRom(romFilepath)

	// Initialize
	data := getData(bytes)
	file := getFile(romFilepath, bytes, data)
	rom := newRom(data, file)

	// Parse it, using the selectedParser
	rom.Content = rom.Parser.Parse(bytes)

	if rom.Parser == nil {
		fmt.Println("Error: Couldn't parse the ROM (unknown format)")
		os.Exit(1)
	}

	return rom
}

func getFile(romFilepath string, bytes []byte, data ContentType) Properties {
	file := Properties{
		Path: romFilepath,
		Name: filepath.Base(romFilepath),
		Size: len(bytes),
		Data: data,
	}
	return file
}

func getData(bytes []byte) ContentType {
	data := ContentType{
		Binary: calc.BinaryDump(bytes),
		Hex:    calc.Hexdump(bytes),
		Bytes:  bytes,
		ASCII:  calc.AsciiDump(calc.Hexdump(bytes)),
	}
	return data
}

func isINES2(romDump ContentType) bool {
	return calc.IsBitSet(romDump.Bytes[7], 3) && !calc.IsBitSet(romDump.Bytes[7], 2)
}

func isINES(romDump ContentType) bool {
	if romDump.ASCII[0] == "N" && romDump.ASCII[1] == "E" && romDump.ASCII[2] == "S" {
		return romDump.Hex[3] == "1A" // this is the <EOF> equivalent in Hex
	}
	return false
}

func selectParser(romDump ContentType) parser {
	if isINES(romDump) {
		if isINES2(romDump) {
			return format.INES2{}
		} else {
			return format.INES{}
		}
	}
	return nil
}
