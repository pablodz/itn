package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {
	processor := itn.NewLanguageES()
	new_string := processor.Alpha2Digit(
		"uno quince",
		false,
		true,
		3,
	)
	println(new_string)
}
