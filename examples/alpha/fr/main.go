package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {
	itn.SetDebug(true)

	processor, _ := itn.NewLanguage(itn.French)
	new_string := processor.Alpha2Digit("un millième", false, true, 3)
	println(new_string)
	println("-----------------------------------------------------")
	println("un 1000ème")
}
