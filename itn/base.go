package itn

import (
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

type RelaxTuple struct {
	Zero string
	One  string
}

func (lg *Language) NotNumericWord(word string) bool {
	isEmpty := false
	if word == "" {
		isEmpty = true
	}

	isDecimalSep := false
	if word == lg.DecimalSep {
		isDecimalSep = true
	}

	isNotNumber := false
	if _, ok := lg.Numbers[word]; !ok {
		isNotNumber = true
	}

	isNotZero := false
	for _, zero := range lg.Zero {
		if word == zero {
			isNotZero = true
			break
		}
	}

	return isEmpty || isDecimalSep && isNotNumber && isNotZero
}

func (lg *Language) Normalize(word string) string {
	return word
}

func (lg *Language) Ord2Card(word string) string {
	return ""
}

func (lg *Language) NumOrd(digits string, originalWord string) string {
	return ""
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

func (lg *SpanishLanguage) Alpha2Digit(text string, relaxed bool, signed bool, ordinalThreshold int) string {
	segments := WORDSEP.Split(text, -1)
	for i, segment := range segments {
		log.Println("[segment]", i, segment)
	}
	punct := WORDSEP.FindAllString(text, -1)
	for i, p := range punct {
		log.Println("[punct]", i, p)
	}

	if len(punct) < len(segments) {
		punct = append(punct, "")
	}

	segmentAndPuncts := []segmentAndPunct{}
	for i, segment := range segments {
		segmentAndPuncts = append(segmentAndPuncts, segmentAndPunct{segment, punct[i]})
	}

	outSegments := []string{}
	for _, sp := range segmentAndPuncts {
		tokens := strings.Split(sp.segment, " ")
		log.Printf("[sp.segment] %s [len]%d", sp.segment, len(tokens))

		numBuilder := NewWordToDigitParser(lg.Language, relaxed, signed, ordinalThreshold, "")
		lastWord := ""
		inNumber := false
		outTokens := []string{}
		for _, couple := range lookAhead(tokens) {

			log.Printf("[word] %s [next] %s", couple.Word, couple.Ahead)

			if numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead)) {
				log.Printf("condition 1: word %s ahead %s", couple.Word, couple.Ahead)
				inNumber = true
			} else if inNumber {
				log.Printf("condition 2: word %s ahead %s", couple.Word, couple.Ahead)
				log.Printf("numBuilder.value() >>>>>>>> %s", numBuilder.GetValue())
				outTokens = append(outTokens, numBuilder.GetValue())
				log.Printf("relaxed %v signed %v ordinalThreshold %d lastWord %s", relaxed, signed, ordinalThreshold, lastWord)
				numBuilder = NewWordToDigitParser(lg.Language, relaxed, signed, ordinalThreshold, lastWord)
				inNumber = numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead))
				log.Printf("inNumber %v", inNumber)
			}

			if !inNumber {
				log.Printf("condition 3: word %s ahead %s", couple.Word, couple.Ahead)
				outTokens = append(outTokens, couple.Word)
			}
			lastWord = strings.ToLower(couple.Word)
		}
		numBuilder.close()
		if numBuilder.GetValue() != "" {
			log.Printf("numBuilder.value() %s", numBuilder.GetValue())
			outTokens = append(outTokens, numBuilder.GetValue())
		}

		outSegments = append(outSegments, strings.Join(outTokens, " "))
		outSegments = append(outSegments, sp.punct)
	}
	text = strings.Join(outSegments, "")

	return text
}
