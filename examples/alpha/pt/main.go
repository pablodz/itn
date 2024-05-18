package main

import (
	"github.com/pablodz/itn/itn"
)

func main() {
	itn.SetDebug(true)

	processor, _ := itn.NewLanguage(itn.Portuguese)
	new_string := processor.Alpha2Digit("Trezentos e setenta e oito milhões vinte e sete mil trezentos e doze", false, true, 3)
	println(new_string)
	println("-----------------------------------------------------")
	println("378027312")
}
