package itn

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Language struct {
	Multipliers                   map[string]int
	Units                         map[string]int
	STens                         map[string]int
	MTens                         map[string]int
	MTensWSTens                   []string
	Hundred                       map[string]int
	MHundreds                     map[string]int
	Numbers                       map[string]int
	Sign                          map[string]string
	Zero                          []string
	DecimalSep                    string
	DecimalSYM                    string
	AndNums                       []string
	And                           string
	NeverIfAlone                  []string
	Relaxed                       map[string]RelaxTuple
	Simplify_check_coef_appliable bool
}

func NewLanguageES() *Language {

	l := &Language{
		Multipliers: map[string]int{
			"mil":      1000,
			"miles":    1000,
			"millon":   1000000,
			"millón":   1000000,
			"millones": 1000000,
		},
		Units: map[string]int{
			"uno":    1,
			"dos":    2,
			"tres":   3,
			"cuatro": 4,
			"cinco":  5,
			"seis":   6,
			"siete":  7,
			"ocho":   8,
			"nueve":  9,
			"un":     1, // optional
			"una":    1, // optional

		},
		STens: map[string]int{
			"diez":         10,
			"once":         11,
			"doce":         12,
			"trece":        13,
			"catorce":      14,
			"quince":       15,
			"dieciseis":    16,
			"diecisiete":   17,
			"dieciocho":    18,
			"diecinueve":   19,
			"veinte":       20,
			"veintiuno":    21,
			"veintidos":    22,
			"veintitres":   23,
			"veinticuatro": 24,
			"veinticinco":  25,
			"veintiseis":   26,
			"veintisiete":  27,
			"veintiocho":   28,
			"veintinueve":  29,
			"veintitrés":   23, // with accent
			"veintidós":    22, // with accent
		},
		MTens: map[string]int{
			"treinta":   30,
			"cuarenta":  40,
			"cincuenta": 50,
			"sesenta":   60,
			"setenta":   70,
			"ochenta":   80,
			"noventa":   90,
		},
		MTensWSTens: []string{},
		Hundred: map[string]int{
			"cien":          100,
			"ciento":        100,
			"cienta":        100,
			"doscientos":    200,
			"trescientos":   300,
			"cuatrocientos": 400,
			"quinientos":    500,
			"seiscientos":   600,
			"setecientos":   700,
			"ochocientos":   800,
			"novecientos":   900,
			"doscientas":    200, // with feminine
			"trescientas":   300, // with feminine
			"cuatrocientas": 400, // with feminine
			"quinientas":    500, // with feminine
			"seiscientas":   600, // with feminine
			"setecientas":   700, // with feminine
			"ochocientas":   800, // with feminine
			"novecientas":   900, // with feminine
		},
		Sign: map[string]string{
			"mas":   "+",
			"menos": "-",
		},
		Zero: []string{
			"cero",
		},
		DecimalSep: "coma",
		DecimalSYM: ".",
		AndNums: []string{
			"un",
			"uno",
			"una",
			"dos",
			"tres",
			"cuatro",
			"cinco",
			"seis",
			"siete",
			"ocho",
			"nueve",
		},

		And: "y",
		NeverIfAlone: []string{
			"un",
			"uno",
			"una",
		},
		Relaxed: map[string]RelaxTuple{},
	}

	// deep copy from l.multipliers
	l.Numbers = map[string]int{
		"mil":      1000,
		"miles":    1000,
		"millon":   1000000,
		"millón":   1000000,
		"millones": 1000000,
	}

	for k, v := range l.Units {
		l.Numbers[k] = v
	}
	for k, v := range l.STens {
		l.Numbers[k] = v
	}
	for k, v := range l.MTens {
		l.Numbers[k] = v
	}
	for k, v := range l.Hundred {
		l.Numbers[k] = v
	}

	return l
}

type RelaxTuple struct {
	Zero string
	One  string
}

func (lg *Language) Ord2Card(word string) string {
	return ""
}

func (lg *Language) NumOrd(digits string, originalWord string) string {
	if strings.HasSuffix(originalWord, "o") {
		return fmt.Sprintf("%sº", digits)
	}
	return fmt.Sprintf("%sª", digits)
}

func (lg *Language) Normalize(word string) string {
	return word
}

func (lg *Language) NotNumericWord(word string) bool {
	return word == "" || word != lg.DecimalSep && !containsKey(lg.Numbers, word) && !contains(lg.Zero, word)
}

var WORDSEP = regexp.MustCompile(`\s*[\.,;\(\)…\[\]:!\?]+\s*|\n`)

type segmentAndPunct struct {
	segment string
	punct   string
}

type LookAhead struct {
	Word  string
	Ahead string
}

func lookAhead(tokens []string) []LookAhead {
	if len(tokens) == 0 {
		return []LookAhead{}
	}

	lookAheads := []LookAhead{}
	for i := 0; i < len(tokens); i++ {

		nextWord := ""
		if i+1 >= len(tokens) {
			nextWord = ""
		} else {
			nextWord = tokens[i+1]
		}

		lookAheads = append(lookAheads, LookAhead{tokens[i], nextWord})
	}
	// fill the last element with empty next
	lookAheads[len(lookAheads)-1].Ahead = ""

	return lookAheads
}

func (lg Language) Alpha2Digit(text string, relaxed bool, signed bool, ordinalThreshold int) string {
	segments := WORDSEP.Split(text, -1)
	// for i, segment := range segments {
	// 	log.Println("[segment]", i, segment)
	// }
	punct := WORDSEP.FindAllString(text, -1)
	// for i, p := range punct {
	// 	log.Println("[punct]", i, p)
	// }

	if len(punct) < len(segments) {
		punct = append(punct, "")
	}

	segmentAndPuncts := []segmentAndPunct{}
	for i, segment := range segments {
		segmentAndPuncts = append(segmentAndPuncts,
			segmentAndPunct{
				segment,
				punct[i],
			},
		)
	}

	outSegments := []string{}
	for _, sp := range segmentAndPuncts {
		tokens := strings.Split(sp.segment, " ")
		log.Printf("tokens %v", tokens)

		numBuilder := NewWordToDigitParser(lg, relaxed, signed, ordinalThreshold, "")
		lastWord := ""
		inNumber := false
		outTokens := []string{}
		for _, couple := range lookAhead(tokens) {

			log.Printf("✅ [word] %s [ahead] %s", couple.Word, couple.Ahead)

			pushed := numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead))
			if pushed {
				log.Printf("> condition 1: word %s ahead %s", couple.Word, couple.Ahead)
				inNumber = true
			} else if inNumber {
				log.Printf("> condition 2: word %s ahead %s", couple.Word, couple.Ahead)
				outTokens = append(outTokens, numBuilder.GetValue())
				numBuilder = NewWordToDigitParser(lg, relaxed, signed, ordinalThreshold, lastWord)
				inNumber = numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead))
			}

			if !inNumber {
				log.Printf("> condition 3: word %s ahead %s", couple.Word, couple.Ahead)
				outTokens = append(outTokens, couple.Word)
			}

			lastWord = strings.ToLower(couple.Word)

			log.Printf("... lastWord %s, inNumber %t, outTokens %v", lastWord, inNumber, outTokens)

		}

		log.Printf("---")
		numBuilder.close()
		if numBuilder.GetValue() != "" {
			outTokens = append(outTokens, numBuilder.GetValue())
		}

		outSegments = append(outSegments, strings.Join(outTokens, " "))
		outSegments = append(outSegments, sp.punct)

	}
	text = strings.Join(outSegments, "")

	return text
}
