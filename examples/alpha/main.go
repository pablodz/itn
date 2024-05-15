package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {

	itn.SetDebug(true)

	processor := itn.NewLanguageES()
	new_string := processor.Alpha2Digit("uno dos quince", false, true, 3)
	println(new_string)
}
