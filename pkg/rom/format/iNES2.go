package format

import (
	"encoding/binary"
	"fmt"
	"github.com/drpaneas/sengo/pkg/calc"
	"hash/crc32"
	"math"
)

// INES2 is a sub-class of section.Sections
type INES2 struct {
	Sections
}

func (i INES2) Parse(romDump []byte) Sections {
	fmt.Println("Rom Heder Version: iNES 2.0")
	romWithoutHeader := romDump[16:]
	romSizeWithoutHeader := len(romWithoutHeader)
	fmt.Printf("Rom Size (without header): %v bytes\n", romSizeWithoutHeader)
	fmt.Printf("Rom Size (without header) CRC32: %08x\n", crc32.Checksum(romWithoutHeader, crc32.IEEETable))

	header := romDump[:16] // header 16 bytes
	trainerSize := 512     // It is always 512 bytes
	var prgrom, chrrom, trainer, chrram, pgrram, playChoiceInstRom, playChoicePROM []byte

	/*	Trainer Area
		-------------
		Trainer exists if bit 2 of Header byte 6 is set.
		It contains data to be loaded into CPU memory at 0x7000
		It is only used by some games that were modified to run on different hardware from the original cartridges,
		such as early RAM cartridges and emulators, adding some compatibility code into those address ranges.
		Trainer is placed between header and PRG ROM data, so PRG ROM should start in the next avail address

	*/
	hasTrainer := false
	if calc.IsBitSet(header[6], 2) {
		//hasTrainer = true
		low := len(header) // the Trainer Area follows the 16-byte Header and precedes the PRG-ROM area
		high := low + trainerSize
		trainer = romDump[low:high]
		trainerSize = len(trainer)
	}

	if hasTrainer {
		fmt.Println("Trainer: Present")
	} else {
		fmt.Println("Trainer: Not present")
		trainerSize = 0
	}
	fmt.Printf("Trainer CRC32: %08x\n", crc32.Checksum(trainer, crc32.IEEETable))


	/*	PRG-ROM Area
		------------
		The PRG-ROM Area follows the 16-byte Header and the Trainer Area and precedes the CHR-ROM Area.
		Header byte 4 (LSB) and bits 0-3 of Header byte 9 (MSB) together specify its size.
		If the MSB nibble is $0-E, LSB and MSB together simply specify the PRG-ROM size in 16 KiB units:
	*/

	// Size of PRG ROM in 16 KB units
	// The PRG-ROM Area follows the 16-byte Header and the Trainer Area (if exists) and precedes the CHR-ROM Area.
	var sizeOfPrgRom, sizeOfPrgRomInKB, sizeOfPrgRomIn16KB int
	MSNibbleByte9 := calc.ReadLowNibbleByte(header[9])
	if calc.ByteToHex(MSNibbleByte9) == "0F" {
		E := (header[4] & 0b11111100) >> 2
		MM := header[4] & 0b00000011
		sizeOfPrgRom = int(math.Pow(2,float64(E)))*(calc.ByteToInt(MM)*2+1)
		sizeOfPrgRomInKB = sizeOfPrgRom / 1024
		sizeOfPrgRomIn16KB = sizeOfPrgRomInKB / 16
	} else {
		tmp := fmt.Sprintf("%v%v", calc.ByteToHex(MSNibbleByte9), calc.ByteToHex(header[4]))
		sizeOfPrgRomIn16KB = calc.HexToInt(tmp)
		sizeOfPrgRomInKB = sizeOfPrgRomIn16KB * 16
		sizeOfPrgRom = sizeOfPrgRomInKB * 1024
	}

	prgrom = romWithoutHeader[trainerSize:trainerSize+sizeOfPrgRom]
	fmt.Printf("Size of PRG ROM: 16KB x %v (%v KB) (%v bytes)\n", sizeOfPrgRomIn16KB, sizeOfPrgRomInKB, sizeOfPrgRom)
	fmt.Printf("PRG ROM CRC32: %08x\n", crc32.Checksum(prgrom, crc32.IEEETable))

	/*	CHR-ROM Area
		------------
		The CHR-ROM Area, if present, follows the Trainer and PRG-ROM Areas and precedes the Miscellaneous ROM Area.
		Header byte 5 (LSB) and bits 4-7 of Header byte 9 (MSB) specify its size.
		If the MSB nibble is $0-E, LSB and MSB together simply specify the CHR-ROM size in 8 KiB units:
	*/

	var sizeOfChrRom, sizeOfChrRomInKB, sizeOfChrRomIn8KB int
	MSBNibbleByte9 := calc.ReadHighNibbleByte(header[9])
	if calc.ByteToHex(MSBNibbleByte9) == "0F" {
		E := (header[5] & 0b11111100) >> 2
		MM := header[5] & 0b00000011
		sizeOfChrRom = int(math.Pow(2,float64(E)))*(calc.ByteToInt(MM)*2+1)
		sizeOfChrRomInKB = sizeOfChrRom / 1024
		sizeOfChrRomIn8KB = sizeOfChrRomInKB / 8
	} else {
		tmp := fmt.Sprintf("%v%v", calc.ByteToHex(MSBNibbleByte9), calc.ByteToHex(header[5]))
		sizeOfChrRomIn8KB = calc.HexToInt(tmp)
		sizeOfChrRomInKB = sizeOfChrRomIn8KB * 8
		sizeOfChrRom = sizeOfChrRomInKB * 1024
	}

	chrrom = romWithoutHeader[trainerSize+sizeOfPrgRom:trainerSize+sizeOfPrgRom+sizeOfChrRom]
	fmt.Printf("Size of CHR ROM: 8KB x %v (%v KB) (%v bytes)\n", sizeOfChrRomIn8KB, sizeOfChrRomInKB, sizeOfChrRom)
	fmt.Printf("CHR ROM CRC32: %08x\n", crc32.Checksum(chrrom, crc32.IEEETable))

	/* 	Miscellaneous ROM Area
		----------------------
		The Miscellaneous ROM Area, if present, follows the CHR-ROM area and occupies the remainder of the file.
		Its size is not explicitly denoted in the header, and can be deduced by subtracting
		the 16-byte Header, Trainer, PRG-ROM and CHR-ROM Area sizes from the total file size.
		The meaning of this data depends on the console type and mapper type;

		Header byte 14 is used to denote the presence of the Miscellaneous ROM Area and
		the number of ROM chips in case any disambiguation is needed.
	 */

	numberOfMiscRoms := header[14] & 0b00000011
	fmt.Printf("Number of Misc Roms: %v\n", numberOfMiscRoms)
	var miscRom []byte
	miscRomSize := romSizeWithoutHeader - (trainerSize + sizeOfPrgRom + sizeOfChrRom)
	if miscRomSize != 0 {
		start := trainerSize + sizeOfPrgRom + sizeOfChrRom
		miscRom = romWithoutHeader[start:]
	}

	fmt.Printf("Size of Misc Rom: %v bytes\n", miscRomSize)
	fmt.Printf("Misc ROM CRC32: %08x\n", crc32.Checksum(miscRom, crc32.IEEETable))

	mapper1 := calc.ReadHighNibbleByte(header[6])   // Lower bits of mapper
	mapper2 := calc.ReadHighNibbleByte(header[7])	// Upper bits of mapper
	mapper3 := calc.ReadLowNibbleByte(header[8])
	mapper := binary.LittleEndian.Uint16([]byte{calc.MergeNibbles(mapper2, mapper1), mapper3})
	fmt.Printf("Mapper Number is: %v\n", mapper)
	// SubMapper number
	subMapper := calc.ReadHighNibbleByte(header[8])
	fmt.Printf("Submapper Number is: %d\n", subMapper)


	// Mirroring
	// TODO: Check the mapper here
	// Header Byte 6 bit 0 is relevant only if the mapper does not allow the mirroring type to be switched.
	// Otherwise, it must be ignored and should be set to zero.
	mirroring := "Ignored"
	if calc.IsBitSet(header[6], 3) {
		fmt.Println("Four-screen VRAM: Yes") // Ignore mirroring control and the mirroring bit
	} else {
		fmt.Println("Four-screen VRAM: No") // Do NOT ignore mirroring
		// Mirroring Typex
		if calc.IsBitSet(header[6], 0) {
			mirroring = "vertical"
		} else {
			mirroring = "horizontal or mapper controlled"
		}
		fmt.Printf("Mirroring: %s\n", mirroring)
	}

	if calc.IsBitSet(header[6], 1) {
		fmt.Println("Battery Backup or other non-volatile memory: Present")
		start := calc.HexToInt("0x6000")
		finish := calc.HexToInt("0x7FFF")
		pgrram = romDump[start:finish]
	} else {
		fmt.Println("Battery backup or other non-volatile memory: Not Present")
	}

	// Console Type
	if !calc.IsBitSet(header[7], 0) && !calc.IsBitSet(header[7], 1) {
		fmt.Println("Console Type: Nintendo Entertainment System/Family Computer")
	}

	if calc.IsBitSet(header[7], 0) && !calc.IsBitSet(header[7], 1) {
		fmt.Println("Console Type: Nintendo Vs System")
	}

	if !calc.IsBitSet(header[7], 0) && calc.IsBitSet(header[7], 1) {
		fmt.Println("Console Type: Nintendo Playchoise 10")
	}

	if calc.IsBitSet(header[7], 0) && calc.IsBitSet(header[7], 1) {
		fmt.Println("Console Type: Extended Console Type")
		consoleTypeByte := header[7] & 0b00000011	// take bit 0 and 1
		consoleType := getExtendedConsoleType(consoleTypeByte)
		fmt.Printf("Extended Console Type: %s\n", consoleType)
		// If it's an extended console then the Vs. System Type has the following PPU and Hardware Type
		vsSystemPPUByte := calc.ReadLowNibbleByte(header[13])
		vsSystemPPU := getVsPPUType(vsSystemPPUByte)
		fmt.Printf("Vs. System PPU Type: %v\n",vsSystemPPU)
		vsSystemTypeByte := calc.ReadHighNibbleByte(header[13])
		vsSystemType := getVsSystemType(vsSystemTypeByte)
		fmt.Printf("Vs. System Type: %v\n", vsSystemType)
	}

	prgramSize := calc.ReadLowNibbleByte(header[10])
	if prgramSize != 0 {
		fmt.Printf("PRM RAM Size: %v bytes\n", 64<<prgramSize)
	}

	prgnvram := calc.ReadHighNibbleByte(header[10])
	if prgnvram != 0 {
		fmt.Printf("PRG-NVRAM/EEPROM Size: %v bytes\n", 64<<prgnvram)
	}

	chrramSize := calc.ReadLowNibbleByte(header[11])
	if chrramSize != 0 {
		fmt.Printf("CHR RAM Size: %v bytes\n", 64<<chrramSize)
	}

	chrnvramSize := calc.ReadHighNibbleByte(header[11])
	if chrnvramSize != 0 {
		fmt.Printf("CHR NVRAM Size: %v bytes\n", 64<<chrnvramSize)
	}

	cpuPPUTiming := calc.ByteToInt(header[12] & 0b0000011)
	if cpuPPUTiming == 0 {
		fmt.Println("CPU/PPU timing mode: RP2C02 (\"NTSC NES\")")
		fmt.Println("Region: North America, Japan, South Korea, Taiwan")
	} else if cpuPPUTiming == 1 {
		fmt.Println("CPU/PPU timing mode: RP2C07 (\"Licensed PAL NES\")")
		fmt.Println("Region: Western Europe, Australia")
	} else if cpuPPUTiming == 2 {
		fmt.Println("CPU/PPU timing mode: Multiple-region")
		fmt.Println("Region: Identical ROM content in both NTSC and PAL countries.")
	} else if cpuPPUTiming == 3 {
		fmt.Println("CPU/PPU timing mode: UMC 6527P (\"Dendy\")")
		fmt.Println("Region: Eastern Europe, Russia, Mainland China, India, Africa")
	} else {
		fmt.Println("CPU/PPU timing mode: Unknown")
	}

	defaultExpansionDeviceByte := header[15] & 0b00111111
	defaultExpansionDevice := getDefaultExpansionDevice(defaultExpansionDeviceByte)
	fmt.Printf("Default Expansion Device: %s\n",defaultExpansionDevice)

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

func getExtendedConsoleType(consoleTypeByte uint8) string {
	switch consoleTypeByte {
	case 3:
		return "Regular Famiclone, but with CPU that supports Decimal Mode (e.g. Bit Corporation Creator)"
	case 4:
		return "V.R. Technology VT01 with monochrome palette"
	case 5:
		return "V.R. Technology VT01 with red/cyan STN palette"
	case 6:
		return "V.R. Technology VT02"
	case 7:
		return "V.R. Technology VT03"
	case 8:
		return "V.R. Technology VT09"
	case 9:
		return "V.R. Technology VT32"
	case 10:
		return "V.R. Technology VT369"
	case 11:
		return "UMC UM6578"
	default:
		return "Unknown/Undefined"
	}
}

func getVsPPUType(vsPPUTypeByte uint8) string {
	switch vsPPUTypeByte {
	case 0:
		return "RP2C03B"
	case 1:
		return "RP2C03G"
	case 2:
		return "RP2C04-0001"
	case 3:
		return "RP2C04-0002"
	case 4:
		return "RP2C04-0003"
	case 5:
		return "RP2C04-0004"
	case 6:
		return "RC2C03B"
	case 7:
		return "RC2C03C"
	case 8:
		return "RC2C05-01 ($2002 AND $?? =$1B)"
	case 9:
		return "RC2C05-02 ($2002 AND $3F =$3D)"
	case 10:
		return "RC2C05-03 ($2002 AND $1F =$1C)"
	case 11:
		return "RC2C05-04 ($2002 AND $1F =$1B)"
	case 12:
		return "RC2C05-05 ($2002 AND $1F =unknown)"
	default:
		return "Unknown/Undefined"
	}
}

func getVsSystemType(vsSystemTypeByte uint8) string {
	switch vsSystemTypeByte {
	case 0:
		return "Vs. Unisystem (normal)"
	case 1:
		return "Vs. Unisystem (RBI Baseball protection)"
	case 2:
		return "Vs. Unisystem (TKO Boxing protection)"
	case 3:
		return "Vs. Unisystem (Super Xevious protection)"
	case 4:
		return "Vs. Unisystem (Vs. Ice Climber Japan protection)"
	case 5:
		return "Vs. Dual System (normal)"
	case 6:
		return "Vs. Dual System (Raid on Bungeling Bay protection)"
	default:
		return "Unknown/Undefined"
	}
}

func getDefaultExpansionDevice(defaultExpansionDeviceByte uint8) string {
	switch defaultExpansionDeviceByte {
	case 0:
		return "Unspecified"
	case 1:
		return "Standard NES/Famicom controllers"
	case 2:
		return "NES Four Score/Satellite with two additional standard controllers"
	case 3:
		return "Famicom Four Players Adapter with two additional standard controllers"
	case 4:
		return "Vs. System"
	case 5:
		return "Vs. System with reversed inputs"
	case 6:
		return "Vs. Pinball (Japan)"
	case 7:
		return "Vs. Zapper"
	case 8:
		return "Zapper ($4017)"
	case 9:
		return "Two Zappers"
	case 10:
		return "Bandai Hyper Shot Lightgun"
	case 11:
		return "Power Pad Side A"
	case 12:
		return "Power Pad Side B"
	case 13:
		return "Family Trainer Side A"
	case 14:
		return "Family Trainer Side B"
	case 15:
		return "Arkanoid Vaus Controller (NES)"
	case 16:
		return "Arkanoid Vaus Controller (Famicom)"
	case 17:
		return "Two Vaus Controllers plus Famicom Data Recorder"
	case 18:
		return "Konami Hyper Shot Controller"
	case 19:
		return "Coconuts Pachinko Controller"
	case 20:
		return "Exciting Boxing Punching Bag (Blowup Doll)"
	case 21:
		return "Jissen Mahjong Controller"
	case 22:
		return "Party Tap"
	case 23:
		return "Oeka Kids Tablet"
	case 24:
		return "Sunsoft Barcode Battler"
	case 25:
		return "Miracle Piano Keyboard"
	case 26:
		return "Pokkun Moguraa (Whack-a-Mole Mat and Mallet)"
	case 27:
		return "Top Rider (Inflatable Bicycle)"
	case 28:
		return "Double-Fisted (Requires or allows use of two controllers by one player)"
	case 29:
		return "Famicom 3D System"
	case 30:
		return "Doremikko Keyboard"
	case 31:
		return "R.O.B. Gyro Set"
	case 32:
		return "Famicom Data Recorder (don't emulate keyboard)"
	case 33:
		return "ASCII Turbo File"
	case 34:
		return "IGS Storage Battle Box"
	case 35:
		return "Family BASIC Keyboard plus Famicom Data Recorder"
	case 36:
		return "Dongda PEC-586 Keyboard"
	case 37:
		return "Bit Corp. Bit-79 Keyboard"
	case 38:
		return "Subor Keyboard"
	case 39:
		return "Subor Keyboard plus mouse (3x8-bit protocol)"
	case 40:
		return "Subor Keyboard plus mouse (24-bit protocol)"
	case 41:
		return "SNES Mouse ($4017.d0)"
	case 42:
		return "Multicart"
	case 43:
		return "Two SNES controllers replacing the two standard NES controllers"
	case 44:
		return "RacerMate Bicycle"
	case 45:
		return "U-Force"
	case 46:
		return "R.O.B. Stack-Up"
	case 47:
		return "City Patrolman Lightgun"
	case 48:
		return "Sharp C1 Cassette Interface"
	case 49:
		return "Standard Controller with swapped Left-Right/Up-Down/B-A"
	case 50:
		return "Excalibor Sudoku Pad"
	case 51:
		return "ABL Pinball"
	case 52:
		return "Golden Nugget Casino extra buttons"
	default:
		return "Unknown/Undefined"
	}
}
