package main

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/rom"
	"github.com/drpaneas/sengo/pkg/utils"
)

func main() {
	romFilepath := utils.GetRomFilepathFromUser()
	r := rom.Open(romFilepath)
	r.Data = r.Parse(r.File.Data.Bytes)

	fmt.Println("ROM Properties:", r.File.Path)
	fmt.Printf("ROM Size: %d bytes\n", r.File.Size)
	fmt.Printf("Rom Name: %s\n", r.Name)

}


