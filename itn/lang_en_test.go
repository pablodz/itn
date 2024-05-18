package itn

import (
	"testing"
)

func TestAlpha2DigitEN(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "twenty-five cows, twelve chickens and one hundred twenty five kg of potatoes.",
			output: "25 cows, 12 chickens and 125 kg of potatoes.",
		},
		{
			input:  "one thousand two hundred sixty-six dollars.",
			output: "1266 dollars.",
		},
		{
			input:  "one two three four twenty fifteen",
			output: "1 2 3 4 20 15",
		},
		{
			input:  "twenty-one, thirty-one.",
			output: "21, 31.",
		},
		{
			input:  "one two three four twenty five.",
			output: "1 2 3 4 25.",
		},
		{
			input:  "one two three four twenty, five.",
			output: "1 2 3 4 20, 5.",
		},
		{
			input:  "thirty-four = thirty four",
			output: "34 = 34",
		},
		{
			input:  "forty five hundred thirty eight dollars and eighteen cents",
			output: "4538 dollars and 18 cents",
		},
		{
			input:  "plus thirty-three nine sixty zero six twelve twenty-one",
			output: "+33 9 60 06 12 21",
		},
		{
			input:  "plus thirty-three nine sixty o six twelve twenty-one",
			output: "+33 9 60 06 12 21",
		},
		{
			input:  "zero nine sixty zero six twelve twenty-one",
			output: "09 60 06 12 21",
		},
		{
			input:  "o nine sixty o six twelve twenty-one",
			output: "09 60 06 12 21",
		},
		{
			input:  "My name is o s c a r.",
			output: "My name is o s c a r.",
		},
		{
			input:  "fifty sixty thirty and eleven",
			output: "50 60 30 and 11",
		},
		{
			input:  "thirteen thousand zero ninety",
			output: "13000 090",
		},
		{
			input:  "thirteen thousand o ninety",
			output: "13000 090",
		},
		{
			input:  "zero",
			output: "0",
		},
		{
			input:  "zero love",
			output: "0 love",
		},
		{
			input:  "Fifth third second twenty-first hundredth one thousand two hundred thirtieth twenty-fifth thirty-eighth forty-ninth.",
			output: "5th third second 21st 100th 1230th 25th 38th 49th.",
		},
		{
			input:  "first, second, third, fourth, fifth, sixth, seventh, eighth, ninth, tenth.",
			output: "first, second, third, 4th, 5th, 6th, 7th, 8th, 9th, 10th.",
		},
		{
			input:  "twenty second position at the twenty first event lost after the first second",
			output: "22nd position at the 21st event lost after the first second",
		},
		{
			input:  "twelve point ninety-nine, one hundred twenty point zero five, one hundred twenty point o five, one point two hundred thirty-six.",
			output: "12.99, 120.05, 120.05, 1.236.",
		},
		{
			input:  "point fifteen",
			output: "0.15",
		},
		{
			input:  "The average density is zero point five",
			output: "The average density is 0.5",
		},
		{
			input:  "This is the one I'm looking for. One moment please! Twenty one cats. One two three four!",
			output: "This is the one I'm looking for. One moment please! 21 cats. 1 2 3 4!",
		},
		{
			input:  "No one is innocent. Another one bites the dust.",
			output: "No one is innocent. Another one bites the dust.",
		},
		{
			input:  "one cannot know",
			output: "one cannot know",
		},
		{
			input:  "the sixth one",
			output: "the 6th one",
		},
		{
			input:  "No one. Another one. One one. Twenty one",
			output: "No one. Another one. 1 1. 21",
		},
		{
			input:  "One second please! twenty second is parsed as twenty-second and is different from twenty seconds.",
			output: "One second please! 22nd is parsed as 22nd and is different from 20 seconds.",
		},
		{
			input:  "FIFTEEN ONE TEN ONE",
			output: "15 1 10 1",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(English)
		new_string := processor.Alpha2Digit(tt.input, false, true, 3)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}

func TestAlpha2DigitENSpecialConfig(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "first, second, third, fourth, fifth, sixth, seventh, eighth, ninth, tenth.",
			output: "1st, 2nd, 3rd, 4th, 5th, 6th, 7th, 8th, 9th, 10th.",
		},
	}

	for _, tt := range tests {
		processor, _ := NewLanguage(English)
		new_string := processor.Alpha2Digit(tt.input, false, true, 0)
		if new_string != tt.output {
			t.Errorf("❌ Expected <%s>, got <%s>", tt.output, new_string)
		} else {
			t.Logf("✅ Expected <%s>, got <%s>", tt.output, new_string)
		}
	}
}
