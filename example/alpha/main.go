package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {
	processor := itn.NewSpanishLanguage()
	new_string := processor.Alpha2Digit(
		"uno dos tres cuatro siete siete noventa 89",
		false,
		true,
		3,
	)
	println(new_string)
}
