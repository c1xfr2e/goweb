package render

import "github.com/bluecover/lm/server/codec"

var (
	encoder   codec.Encoder
	printData bool
)

func SetEncoder(enc codec.Encoder) {
	encoder = enc
}

func PrintData(b ...bool) {
	if len(b) == 0 {
		printData = true
	} else {
		printData = b[0]
	}
}
