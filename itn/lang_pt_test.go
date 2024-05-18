package itn

import (
	"testing"
)

func TestAlpha2DigitPT(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "um vírgula um",
			output: "1,1",
		},
		{
			input:  "um vírgula quatrocentos e um",
			output: "1,401",
		},
		{
			input:  "vinte cinco vacas, doze galinhas e cento vinte e cinco kg de batatas.",
			output: "25 vacas, 12 galinhas e 125 kg de batatas.",
		},
		{
			input:  "mil duzentos sessenta e seis dólares.",
			output: "1266 dólares.",
		},
		{
			input:  "um dois três quatro vinte quinze",
			output: "1 2 3 4 20 15",
		},
		{
			input:  "vinte e um, trinta e um.",
			output: "21, 31.",
		},
		{
			input:  "mais trinta e três nove sessenta zero seis doze vinte e um",
			output: "+33 9 60 06 12 21",
		},
		{
			input:  "zero nove sessenta zero seis doze vinte e um",
			output: "09 60 06 12 21",
		},
		{
			input:  "cinquenta sessenta trinta onze",
			output: "50 60 30 11",
		},
		{
			input:  "duzentos e quarenta e quatro",
			output: "244",
		},
		{
			input:  "dois mil e vinte",
			output: "2020",
		},
		{
			input:  "mil novecentos e oitenta e quatro",
			output: "1984",
		},
		{
			input:  "mil e novecentos",
			output: "1900",
		},
		{
			input:  "dois mil cento e vinte cinco",
			output: "2125",
		},
		{
			input:  "Trezentos e setenta e oito milhões vinte e sete mil trezentos e doze",
			output: "378027312",
		},
		{
			input:  "treze mil zero noventa",
			output: "13000 090",
		},
		{
			input:  "zero",
			output: "0",
		},
		{
			input:  "doze vírgula noventa e nove, cento e vinte vírgula zero cinco, um vírgula duzentos e trinta e seis, um vírgula dois três seis.",
			output: "12,99, 120,05, 1,236, 1,2 3 6.",
		},
		{
			input:  "vírgula quinze",
			output: "0,15",
		},
		{
			input:  "Temos mais vinte graus dentro e menos quinze fora.",
			output: "Temos +20 graus dentro e -15 fora.",
		},
		{
			input:  "Um momento por favor! trinta e um gatos. Um dois três quatro!",
			output: "Um momento por favor! 31 gatos. 1 2 3 4!",
		},
		{
			input:  "Nem um. Um um. Trinta e um",
			output: "Nem um. 1 1. 31",
		},
		{
			input:  "Um milhao",
			output: "1000000",
		},
		{
			input:  "Um segundo por favor! Vigésimo segundo é diferente de vinte segundos.",
			output: "Um segundo por favor! 22º é diferente de 20 segundos.",
		},
		{
			input:  "Ordinais: primeiro, quinto, terceiro, vigésima, vigésimo primeiro, centésimo quadragésimo quinto",
			output: "Ordinais: primeiro, 5º, terceiro, 20ª, 21º, 145º",
		},
		{
			input:  "A décima quarta brigada do exército português, juntamento com o nonagésimo sexto regimento britânico, bateu o centésimo vigésimo sétimo regimento de infantaria de Napoleão",
			output: "A 14ª brigada do exército português, juntamento com o 96º regimento britânico, bateu o 127º regimento de infantaria de Napoleão",
		},
		{
			input:  "em mil quinhentos e catorze, ela nasceu",
			output: "em 1514, ela nasceu",
		},
		{
			input:  "tudo aconteceu até mil novecentos e dezesseis",
			output: "tudo aconteceu até 1916",
		},
		{
			input:  "em dezessete de janeiro de mil novecentos e noventa",
			output: "em 17 de janeiro de 1990",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(Portuguese)
		new_string := processor.Alpha2Digit(tt.input, false, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitPTFalse(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "quanto é dezenove menos três? É dezesseis",
			output: "quanto é 19 menos 3? É 16",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(Portuguese)
		new_string := processor.Alpha2Digit(tt.input, true, false, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitPTRelaxed(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "um dois três quatro trinta e cinco.",
			output: "1 2 3 4 35.",
		},
		{
			input:  "um dois três quatro vinte, cinco.",
			output: "1 2 3 4 20, 5.",
		},
		{
			input:  "trinta e quatro = trinta quatro",
			output: "34 = 34",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(Portuguese)
		new_string := processor.Alpha2Digit(tt.input, true, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}
