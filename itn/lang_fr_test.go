package itn

import (
	"testing"
)

func TestAlpha2DigitFR(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "Vingt-cinq vaches, douze poulets et cent vingt-cinq kg de pommes de terre.",
			output: "25 vaches, 12 poulets et 125 kg de pommes de terre.",
		},
		{
			input:  "Mille deux cent soixante-six clous.",
			output: "1266 clous.",
		},
		{
			input:  "Mille deux cents soixante-six clous.",
			output: "1266 clous.",
		},
		{
			input:  "Nonante-cinq = quatre-vingt-quinze",
			output: "95 = 95",
		},
		{
			input:  "Nonante cinq = quatre-vingt quinze",
			output: "95 = 95",
		},
		{
			input:  "un deux trois quatre vingt quinze",
			output: "1 2 3 4 20 15",
		},
		{
			input:  "Vingt et un, trente et un.",
			output: "21, 31.",
		},
		{
			input:  "trente-quatre = trente quatre",
			output: "34 = 34",
		},
		{
			input:  "plus trente-trois neuf soixante zéro six douze vingt et un",
			output: "+33 9 60 06 12 21",
		},
		{
			input:  "zéro neuf soixante zéro six douze vingt et un",
			output: "09 60 06 12 21",
		},
		{
			input:  "cinquante soixante trente et onze",
			output: "50 60 30 11",
		},
		{
			input:  "treize mille zéro quatre-vingt-dix",
			output: "13000 090",
		},
		{
			input:  "treize mille zéro quatre-vingts",
			output: "13000 080",
		},
		{
			input:  "zéro",
			output: "0",
		},
		{
			input:  "a a un trois sept trois trois sept cinq quatre zéro c c",
			output: "a a 1 3 7 3 3 7 5 4 0 c c",
		},
		{
			input:  "sept un zéro",
			output: "7 1 0",
		},
		{
			input:  "Cinquième premier second troisième vingt et unième centième mille deux cent trentième.",
			output: "5ème premier second troisième 21ème 100ème 1230ème.",
		},
		{
			input:  "un millième",
			output: "un 1000ème",
		},
		{
			input:  "un millionième",
			output: "un 1000000ème",
		},
		{
			input:  "Douze virgule quatre-vingt dix-neuf, cent vingt virgule zéro cinq, un virgule deux cent trente six.",
			output: "12,99, 120,05, 1,236.",
		},
		{
			input:  "la densité moyenne est de zéro virgule cinq.",
			output: "la densité moyenne est de 0,5.",
		},
		{
			input:  "Il fait plus vingt degrés à l'intérieur et moins quinze à l'extérieur.",
			output: "Il fait +20 degrés à l'intérieur et -15 à l'extérieur.",
		},
		{
			input:  "J'en ai vu au moins trois dans le jardin, et non plus deux.",
			output: "J'en ai vu au moins 3 dans le jardin, et non plus 2.",
		},
		{
			input:  "Ne pas confondre un article ou un nom avec un chiffre et inversement : les uns et les autres ; une suite de chiffres : un, deux, trois !",
			output: "Ne pas confondre un article ou un nom avec un chiffre et inversement : les uns et les autres ; une suite de chiffres : 1, 2, 3 !",
		},
		{
			input:  "Je n'en veux qu'un. J'annonce: le un",
			output: "Je n'en veux qu'un. J'annonce: le un",
		},
		{
			input:  "dix + deux\n= douze",
			output: "10 + 2\n= 12",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(French)
		new_string := processor.Alpha2Digit(tt.input, false, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitFRRelaxed(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "un deux trois quatre vingt quinze.",
			output: "1 2 3 95.",
		},
		{
			input:  "Quatre, vingt, quinze, quatre-vingts.",
			output: "4, 20, 15, 80.",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(French)
		new_string := processor.Alpha2Digit(tt.input, true, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitFROrdinal0(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "Cinquième premier second troisième vingt et unième centième mille deux cent trentième.",
			output: "5ème 1er 2nd 3ème 21ème 100ème 1230ème.",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(French)
		new_string := processor.Alpha2Digit(tt.input, false, true, 0)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitFRSignedTrue(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "J'en ai vu au moins trois dans le jardin, et non plus deux.",
			output: "J'en ai vu au moins 3 dans le jardin, et non plus 2.",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(French)
		new_string := processor.Alpha2Digit(tt.input, false, false, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitFRSignedFalse(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "J'en ai vu au moins trois dans le jardin, et non plus deux.",
			output: "J'en ai vu au moins 3 dans le jardin, et non plus 2.",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(French)
		new_string := processor.Alpha2Digit(tt.input, false, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}
