package format

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/calc"
	"hash/crc32"
	"os"
)

// INES is a sub-class of Sections
type INES struct {
	Sections
}

func (i INES) Parse(romDump []byte) Sections {
	header := romDump[:16] // header 16 bytes
	var prgrom, chrrom, trainer, chrram, pgrram, playChoiceInstRom, playChoicePROM []byte

	// Size of Program ROM (in 16 KB units)
	sizeOfPrgRom := calc.ByteToInt(header[4])

	// Size of Character ROM (in 8 KB units)
	sizeOfChrRom := calc.ByteToInt(header[5])
	// TODO: Value 0 means the board uses CHR RAM
	if sizeOfChrRom == 0 {
		fmt.Println("The board uses CHR RAM")
		chrram = chrrom
	}

	sizeOfCombinedRoms := sizeOfPrgRom*16*1024 + sizeOfChrRom*8*1024
	sizeRomWithoutHeader := len(romDump[16:])
	romWithoutHeader := romDump[16:]

	if sizeRomWithoutHeader == sizeOfCombinedRoms {
		// doesn't have any other sections
		prgrom = romWithoutHeader[:sizeOfPrgRom*16*1024]
		chrrom = romWithoutHeader[sizeOfPrgRom*16*1024:]
	}

	// Flags 6
	// ------------------------------------------------------------------------------------------------------
	// Mirroring
	mirroring := "Ignored"
	if calc.IsBitSet(header[6], 3) {
		fmt.Println("Four-screen VRAM: Yes") // Ignore mirroring control and the mirroring bit
	} else {
		fmt.Println("Four-screen VRAM: No")  // Do NOT ignore mirroring
		if calc.IsBitSet(header[6], 0) {
			mirroring = "vertical"
		} else {
			mirroring = "horizontal or mapper controlled"
		}
		fmt.Printf("Mirroring: %s\n", mirroring)
	}

	// Battery or any other non-volatile memory (RAM)
	if calc.IsBitSet(header[6], 1) {
		fmt.Println("Battery backup or other non-volatile memory: Yes")
		start := calc.HexToInt("0x6000")
		finish := calc.HexToInt("0x7FFF")
		pgrram = romDump[start:finish]
	} else {
		fmt.Println("Battery backup or other non-volatile memory: No")

	}

	// Trainer exists if bit 2 of Header byte 6 is set.
	// It contains data to be loaded into CPU memory at 0x7000
	// It is only used by some games that were modified to run on different hardware from the original cartridges,
	// such as early RAM cartridges and emulators, adding some compatibility code into those address ranges.
	if calc.IsBitSet(header[6], 2) {
		size := 512
		fmt.Println("Trainer: 512-byte trainer found (placed between header and PRG ROM data")
		low := 16 // the Trainer Area follows the 16-byte Header and precedes the PRG-ROM area
		high := low+size // It is always 512 bytes in size if present
		trainer = romDump[low:high]
	} else {
		fmt.Println("Trainer: Not present")
	}

	// Mapper Lower Nibble
	lowerNibbleMapper := calc.ReadHighNibbleByte(header[6])

	// Flags 7
	// ------------------------------------------------------------------------------------------------------
	// Upper nibble of mapper number
	upperNibbleMapper := calc.ReadHighNibbleByte(header[7])
	mapperNumber := calc.ByteToInt(calc.MergeNibbles(upperNibbleMapper, lowerNibbleMapper))
	fmt.Printf("Mapper Number is: %d\n", mapperNumber)

	// VS
	if calc.IsBitSet(header[7], 0) {
		fmt.Println("Console Type: Nintendo Vs System")
	}

	// PlayChoice
	if calc.IsBitSet(header[7], 1) {
		fmt.Println("Console Type: Nintendo Playchoise 10")
	}

	// Reserved, must be zeroes!
	if calc.IsBitSet(header[7], 2){
		fmt.Println("This bit 2 of Byte 7, should be zero, but it's not!")
	}

	// Otherwise, this is an iNES 2.0 Format
	if calc.IsBitSet(header[7], 3){
		fmt.Println("This bit 3 of Byte 7, should be zero, but it's not!")
	}

	// Flags 8
	// ------------------------------------------------------------------------------------------------------
	// Size of PRG RAM in 8 KB units
	// The PRG RAM Size value (stored in byte 8) was recently added to the official specification;
	// as such, virtually no ROM images in circulation make use of it.
	sizeOfPrgRam := calc.ByteToInt(header[8])
	// Value 0 infers 8 KB for compatibility
	if sizeOfPrgRam == 0 {
		sizeOfPrgRam = 1
	}
	fmt.Printf("Size of PRG RAM: 8KB x %v (%v KB) (%v bytes)\n", sizeOfPrgRam, sizeOfPrgRam*8, sizeOfPrgRam*8*1024)

	// Flags 9
	// ------------------------------------------------------------------------------------------------------
	// Though in the official specification
	// very few emulators honor this bit as virtually no ROM images in circulation make use of it.
	if calc.IsBitSet(header[9], 0) {
		fmt.Println("TV system: PAL")
	} else {
		fmt.Println("TV system: NTSC")
	}

	// Flags 10
	// ------------------------------------------------------------------------------------------------------
	// This byte is not part of the official specification, and relatively few emulators honor it.
	if !calc.IsBitSet(header[10], 0) && !calc.IsBitSet(header[10], 1) {
		fmt.Println("TV system: NTSC")
	}
	if !calc.IsBitSet(header[10], 0) && calc.IsBitSet(header[10], 1) {
		fmt.Println("TV system: PAL")
	}
	if calc.IsBitSet(header[10], 0) && !calc.IsBitSet(header[10], 1) || !calc.IsBitSet(header[10], 0) && calc.IsBitSet(header[10], 1) {
		fmt.Println("Dual compatible")
	}

	// PRG RAM present or not ($6000-$7FFF)
	if calc.IsBitSet(header[10], 4) {
		fmt.Println("PRG RAM is not present")
	} else {
		fmt.Println("PRG RAM is present")
	}

	// Board bus conflicts
	if calc.IsBitSet(header[10], 5) {
		fmt.Println("Board has bus conflicts")
	} else {
		fmt.Println("Board has no bus conflicts.")
	}

	// 8 - 15 bytes are not used, and should be 0
	// ------------------------------------------------------------------------------------------------------
	isbytes8to15NotSet := true
	for i:=8; i<=15; i++ {
		if header[i] != 0 {
			isbytes8to15NotSet = false
			fmt.Printf("Byte %v should be zero, but it's %v instead.\n", i, header[i])
		}
	}
	if isbytes8to15NotSet {
		fmt.Println("8 - 15 bytes are not used, and they are 0.")
	}

	fmt.Println("iNES Header")
	fmt.Printf("Size of PRG ROM: 16KB x %v (%v KB) (%v bytes)\n", sizeOfPrgRom, sizeOfPrgRom*16, sizeOfPrgRom*16*1024)
	fmt.Printf("Size of CHR ROM: 8KB x %v (%v KB) (%v bytes)\n", sizeOfChrRom, sizeOfChrRom*8, sizeOfChrRom*8*1024)
	fmt.Printf("Size of Roms combined: %v KB (%v bytes)\n", sizeOfCombinedRoms/1024, sizeOfCombinedRoms)
	fmt.Printf("Rom size without header: %v KB (%v bytes)\n", sizeRomWithoutHeader/1024, sizeRomWithoutHeader)

	fmt.Printf("PRGROM CRC32 %08x\n", crc32.Checksum(prgrom, crc32.IEEETable))
	fmt.Printf("CHRROM CRC32 %08x\n", crc32.Checksum(chrrom, crc32.IEEETable))
	fmt.Printf("Rom wihtout header CRC32 %08x\n", crc32.Checksum(romWithoutHeader, crc32.IEEETable))

	fmt.Println("----------")
	fmt.Printf("Header: (%v-%v)\n", 0, len(header)-1)
	fmt.Printf("PRGROM: (%v-%v) or (%v-%v)\n", len(header)-1+1, len(header)-1+1-1+len(prgrom), calc.IntToHex(len(header)-1+1), calc.IntToHex(len(header)-1+1-1+len(prgrom)))
	fmt.Printf("CHRROM: (%v-%v) or (%v-%v)\n", len(header)-1+1-1+len(prgrom)+1, len(header)-1+1-1+len(prgrom)+1-1+len(chrrom),calc.IntToHex(len(header)-1+1-1+len(prgrom)+1), calc.IntToHex(len(header)-1+1-1+len(prgrom)+1-1+len(chrrom)))
	fmt.Printf("PRGRAM: (%v-%v)\n", 24576, 24576-1+sizeOfPrgRam*8*1024)

	if err := os.WriteFile("header.bin",header,0644); err != nil {
		fmt.Println("Cannot save header.bin")
	}

	if err := os.WriteFile("PRGROM.bin",prgrom,0644); err != nil {
		fmt.Println("Cannot save PRGROM.bin")
	}

	if err := os.WriteFile("CHRROM.bin",chrrom,0644); err != nil {
		fmt.Println("Cannot save CHRROM.bin")
	}

	if err := os.WriteFile("PRGRAM.bin",pgrram,0644); err != nil {
		fmt.Println("Cannot save PRGRAM.bin")
	}


	return Sections{
		Header:            header,
		Trainer:           trainer,
		PRGROM:            prgrom,
		CHRROM:            chrrom,
		PRGRAM:            pgrram,
		CHRRAM:            chrram,
		PlayChoiceInstRom: playChoiceInstRom,
		PlayChoicePROM:    playChoicePROM,
	}
}
