package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Exists reports whether the named file or directory exists.
func exists(path string, isDir bool) bool {
	if path == "" {
		fmt.Println("Path is empty")
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			return false
		}
	}

	return isDir == info.IsDir()
}

// fileExists reports whether the provided file exists.
func fileExists(path string) bool {
	return exists(path, false)
}

func getFilepath() (string, bool) {
	// return "/Users/drpaneas/Downloads/NintendoMultiRomCollectionByGhostware/Super Mario World 9 (Unl) [!].nes", true
	if len(os.Args) == 2 {
		romFile := os.Args[1]
		// Just a temp

		if !fileExists(romFile) {
			fmt.Printf("Sorry, '%s' file not found! You need to provide a valid filepath to your ROM.\n", romFile)
			os.Exit(1)
		}
		return os.Args[1], true
	}
	fmt.Printf("Sorry, you need to provide a file first. Please try: `%s rom.nes`)\n", os.Args[0])
	return "", false
}

func GetFileInBytes(data *bufio.Reader) []byte {
	var fileInBytes []byte
	for {
		buf := make([]byte, 1)
		_, err := data.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Couldn't get file in bytes. Please file a bug")
				fmt.Println(err)
			}
			break // EOF
		}
		fileInBytes = append(fileInBytes, buf[0])
	}
	return fileInBytes
}

func ReadRom(file string) []byte {
	f, err := os.Open(file)

	if err != nil {
		println("Error opening the file. Please open a bug.")
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing the file. Please open a bug.")
			os.Exit(1)
		}
	}(f)

	return GetFileInBytes(bufio.NewReader(f))
}

func GetRomFilepathFromUser() string {
	romFilepath, ok := getFilepath()
	if !ok {
		os.Exit(1)
	}
	return romFilepath
}