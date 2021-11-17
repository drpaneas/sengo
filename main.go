package main

import (
	"fmt"
	"github.com/drpaneas/sengo/pkg/calc"
	"github.com/drpaneas/sengo/pkg/rom"
	"github.com/drpaneas/sengo/pkg/utils"
	"image"
	"image/png"
	"os"
)

func printTile(t [64]uint8) {
	c := 0
	for _, v := range t {
		fmt.Printf("%v", v)
		c++
		if c == 8 {
			fmt.Println()
			c = 0
		}
	}
}

func fixTile(t [64]uint8) [64]uint8 {
	fixed, leftPart, rightPart := t, t, t
	for start := 0; start <= start+7; start = start + 8 {
		end := start + 7
		if end >= 64 {
			break
		}
		c := 0
		for i := start; i <= end; i++ {
			finish := end - c
			c++
			if c > 3 {
				fixed[i] = rightPart[finish]
			} else {
				fixed[i] = leftPart[finish]
			}
		}
	}

	return fixed
}

func merge(graphics uint8, colors uint8) [8]uint8 {
	// Split the graphics byte into individual bits
	var graphicsBits [8]uint8
	for i := 0; i < 8; i++ {
		if calc.IsBitSet(graphics, i) {
			graphicsBits[i] = 1
		}
	}

	// Split the color byte into individual bits
	var colorBits [8]uint8
	for i := 0; i < 8; i++ {
		if calc.IsBitSet(colors, i) {
			colorBits[i] = 1
		}
	}

	// Combine them
	var combineBits [8]uint8
	for i := 0; i < 8; i++ {
		highBit := colorBits[i] << 1
		lowBit := graphicsBits[i]
		combineBits[i] = highBit | lowBit
	}

	return combineBits
}

func main() {
	romFilepath := utils.GetRomFilepathFromUser()
	r := rom.Open(romFilepath)

	fmt.Println("ROM Properties:", r.File.Path)
	fmt.Printf("ROM Size: %d bytes\n", r.File.Size)
	fmt.Printf("Rom Name: %s\n", r.Name)

	chrrom := r.Content.CHRROM
	var tileCounter,byteCounter,bitCounter int

	// Initialize
	skip := false 		// Skips the next 8 bytes because they are part of the colors
	var data [][64]uint8	// 8x8 resolution, so 64 pixels in total.
	data = append(data, [64]uint8{})	// create the starting [0] tile (empty values inside)

	for k := 0; k+8 < len(chrrom); k++ {
		if skip {
			k = k + 8
			skip = false
		}

		// Create the merged Byte from Graphics and Colors (e.g. [0 0 0 0 0 0 0 0]) and it has length 8
		tileByte := merge(chrrom[k], chrrom[k+8])
		tileByteIndex := 0

		for i := bitCounter; i < (bitCounter + 8); i++ {
			data[tileCounter][i] = tileByte[tileByteIndex]
			tileByteIndex++
		}
		// Debug
		//if tileCounter == 146 {
		//	fmt.Printf("Byte[%v]: %v\tk=%v\n", byteCounter, tileByte, k)
		//}

		// If we reached the 64th iteration, this tile is finished (aka no more pixels/bits). Reset the counter.
		bitCounter += 8
		if bitCounter >= 64 {
			bitCounter = 0
		}

		// If we reached the 8th byte then it's time to reset it as well.
		byteCounter++
		if byteCounter == 8 {
			byteCounter = 0
			tileCounter++
			skip = true

			// Respect CHR ROM array memory boundaries.
			if k+9 < len(chrrom) {
				data = append(data, [64]uint8{})
			}
		}
	}

	// Debug
	//fmt.Println("Tile counter:", tileCounter)
	//for t := 0; t < tileCounter; t++ {
	//	// debug
	//	if t == 146 {
	//		fmt.Printf("Tile %v: %v\n", t, data[t][:])
	//	}
	//}

	// Reorder the bits of every byte (swap low nibble with high nibble)
	// This is because when we read the byte, the computer counts like this:
	// e.g. [1 2 2 0 0 3 1 1 ] <-- [7-bit 6-bit 5-bit 4-bit 3-bit 2-bit 1-bit 0-bit]
	// so the top left corner would be the 7-bit of the 0-Byte.
	// But, when we render an image, the top left pixel should be 0-bit of the 0-Byte.
	// So to display this properly, we should reorder the bits
	// e.g.  [1 1 3 0 0 2 2 1] <-- [1-bit, 2-bit, 3-bit, 4-bit, 5-bit, 6-bit, 7-bit]
	for i := 0; i < tileCounter; i++ {
		data[i] = fixTile(data[i])
	}

	fmt.Printf("There are %v tiles in this CHR ROM\n", tileCounter)

	for tile := 0; tile < tileCounter; tile++ {
		// Create the canvas dimension and the color palette
		img := image.NewGray(image.Rect(0, 0, 8, 8)) // Top Left (0,0) and Top Right (8,8) coordinates

		// Put the tile into the image
		img.Pix = data[tile][:]
		for pixel := 0; pixel < 64; pixel++ {
			// The color data is mapped into a 2-bit greyscale palette
			if img.Pix[pixel] == 0 {
				img.Pix[pixel] = 0 // Black
			} else if img.Pix[pixel] == 1 {
				img.Pix[pixel] = 85 // Light Grey
			} else if img.Pix[pixel] == 2 {
				img.Pix[pixel] = 170 // Light Black
			} else if img.Pix[pixel] == 3 {
				img.Pix[pixel] = 255 // White
			} else {
				fmt.Println("Error: Unknown bit value", img.Pix[pixel])
				os.Exit(1)
			}
		}

		// The graphic assets are split in half into 2 (left & right) banks
		bank := "Left"
		if tile < tileCounter/2 {
			bank = "Right"
		}

		// outputFile is a File type which satisfies Writer interface
		filename := fmt.Sprintf("Tile%v-Bank%v.png", tile, bank)
		outputFile, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error: Couldn't save the image:", filename)
			os.Exit(1)
		}

		// Encode takes a writer interface and an image interface
		// We pass it the File and the Greyscale
		if err = png.Encode(outputFile, img); err != nil {
			fmt.Printf("Error: Couldn't encode the %v image to greyscale\n", filename)
			os.Exit(1)
		}

		// Don't forget to close files
		if err = outputFile.Close(); err != nil {
			fmt.Println("Error: Cannot close the file (memory leak issue?)")
			os.Exit(1)
		}
	}
}
