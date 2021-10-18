package rom

import "github.com/drpaneas/sengo/pkg/rom/cartridge"

// parser is the strategy interface that includes the Parse() method
// which can be extended into support reading ROM of various formats (e.g. iNES, iNES 2.0, e.t.c.)
type parser interface {
	// Parse reads the bytecode of a ROM and extracts it into multiple sections based upon a defined format.
	// The compatible formats can be found into the rom/format folder (namely iNES and iNES 2.0):
	// 		1. iNES  : https://wiki.nesdev.org/w/index.php/INES
	//		2. iNES 2: https://wiki.nesdev.org/w/index.php?title=NES_2.0
	// If you want to add another format in the future, you can do it by extending the Parse method.
	// To do that, create a new file int rom/format folder and then create a struct representing your new specific
	// parsing strategy. Then put the method Parse to be a member function of this.
	Parse(rom []byte) cartridge.ContentInSection
}
