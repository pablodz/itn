package itn

import (
	"testing"
)

func TestAlpha2DigitES(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "uno coma uno",
			output: "1.1",
		},
		{
			input:  "uno coma cuatrocientos uno",
			output: "1.401",
		},
		{
			input:  "veinticinco vacas, doce gallinas y ciento veinticinco kg de patatas.",
			output: "25 vacas, 12 gallinas y 125 kg de patatas.",
		},
		{
			input:  "Habían trescientos hombres y quinientas mujeres",
			output: "Habían 300 hombres y 500 mujeres",
		},
		{
			input:  "mil doscientos sesenta y seis dolares.",
			output: "1266 dolares.",
		},
		{
			input:  "un dos tres cuatro veinte quince",
			output: "1 2 3 4 20 15",
		},
		{
			input:  "uno cero treinta y siete", // Telephony first
			output: "1 037",                    // Telephony first
		},
		{
			input:  "veintitrés dos cuatro", // Telephony first
			output: "23 2 4",                // Telephony first
		},
		{
			input:  "veintiuno, treinta y uno.",
			output: "21, 31.",
		},
		{
			input:  "un dos tres cuatro treinta cinco.",
			output: "1 2 3 4 35.",
		},
		{
			input:  "un dos tres cuatro veinte, cinco.",
			output: "1 2 3 4 20, 5.",
		},
		{
			input:  "treinta y cuatro = treinta cuatro",
			output: "34 = 34",
		},
		{
			input:  "mas treinta y tres nueve sesenta cero seis doce veintiuno",
			output: "+33 9 60 06 12 21",
		},
		{
			input:  "cero nueve sesenta cero seis doce veintiuno",
			output: "09 60 06 12 21",
		},
		{
			input:  "cincuenta sesenta treinta y once",
			output: "50 60 30 y 11",
		},
		{
			input:  "trece mil cero noventa",
			output: "13000 090",
		},
		{
			input:  "cero",
			output: "0",
		},
		{
			input:  "doce coma noventa y nueve, ciento veinte coma cero cinco, uno coma doscientos treinta y seis, uno coma dos tres seis.",
			output: "12.99, 120.05, 1.236, 1.2 3 6.",
		},
		{
			input:  "coma quince",
			output: "0.15",
		},
		{
			input:  "Tenemos mas veinte grados dentro y menos quince fuera.",
			output: "Tenemos +20 grados dentro y -15 fuera.",
		},
		{
			input:  "Un momento por favor! treinta y un gatos. Uno dos tres cuatro!",
			output: "Un momento por favor! 31 gatos. 1 2 3 4!",
		},
		{
			input:  "Ni uno. Uno uno. Treinta y uno", // Telephony first
			output: "Ni 1. 1 1. 31",                  // Telephony first
		},
		{
			input:  "un millon",
			output: "1000000",
		},
		{
			input:  "un millón",
			output: "1000000",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(Spanish)
		new_string := processor.Alpha2Digit(tt.input, false, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}
