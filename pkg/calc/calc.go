package calc

import (
	"encoding/hex"
	"fmt"
	"github.com/drpaneas/sengo/pkg/utils"
	"os"
	"strconv"
	"strings"
)

func AsciiDump(hex []string) []string {
	var asciiData []string
	for _, data := range hex {
		asciiData = append(asciiData, HexToASCII(data))
	}
	return asciiData
}

func BinaryDump(b []byte) []string {
	var binData []string
	for _, data := range b {
		binData = append(binData, fmt.Sprintf("%08b", data))
	}
	return binData
}

func Hexdump(b []byte) []string {
	var hexData []string
	for _, data := range b {
		hexData = append(hexData, ByteToHex(data))
	}
	return hexData
}

func ByteToHex(b byte) string {
	bs := make([]byte, 1)
	bs[0] = b
	return strings.ToUpper(hex.EncodeToString(bs))
}

func ByteToASCII(b byte) string {
	str := fmt.Sprintf("%x", b)
	tmp := fmt.Sprintf("%#v", HexToASCII(str))
	return utils.FixFormat(tmp)
}

func HexToASCII(str string) string {
	bs, err := hex.DecodeString(str)
	if err != nil {
		fmt.Printf("Error: Failed to disassemble: %v", err)
		os.Exit(1)
	}
	return string(bs)
}

func HexToInt(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

func BinToInt(binary string) int64 {
	output, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		fmt.Println("Could not convert binary to integer")
		fmt.Println(err)
		os.Exit(1)
	}
	return output
}