package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {
	itn.SetDebug(true)

	processor, _ := itn.NewLanguage(itn.English)
	new_string := processor.Alpha2Digit("first, second, third, fourth, fifth, sixth, seventh, eighth, ninth, tenth.", false, true, 3)
	println(new_string)
}
