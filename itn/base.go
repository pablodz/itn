package itn

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Language struct {
	LangCode                      LanguageCode
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
	Simplify_check_coef_appliable bool              // Optional
	RadMap                        map[string]string // Optional
	Composites                    map[string]int    // Optional
	PtOrdinals                    map[string]string // Only for Portuguese
	IrrOrd                        map[string]RelaxTuple
}

type RelaxTuple struct {
	Zero string
	One  string
}

func (lg *Language) Ord2Card(word string) string {
	switch lg.LangCode {
	case French:
		logPrintf(">>>> Ord2Card.0 [word] %s", word)
		if containsKey(lg.IrrOrd, word) {
			logPrintf(">>>> Ord2Card.1 %s", word)
			return lg.IrrOrd[word].Zero
		}

		plurSuff := strings.HasSuffix(word, "ièmes")
		singSuff := strings.HasSuffix(word, "ième")
		if !(plurSuff || singSuff) {
			logPrintf(">>>> Ord2Card.2 %s", word)
			return ""
		}

		source := ""
		runeCount := utf8.RuneCountInString(word)
		if plurSuff {
			source = string([]rune(word)[:runeCount-5])
			logPrintf(">>>> Ord2Card.3.0.1 %s", source)
		} else {
			source = string([]rune(word)[:runeCount-4])
			logPrintf(">>>> Ord2Card.3.0.2 %s", source)
		}

		if source == "cinqu" {
			source = "cinq"
		} else if source == "neuv" {
			source = "neuf"
		} else if !containsKey(lg.Numbers, source) {
			source = source + "e"
			logPrintf(">>>> Ord2Card.3.1 %s", source)
			if !containsKey(lg.Numbers, source) {
				logPrintf(">>>> Ord2Card.3.2 %s", source)
				return ""
			}
		}
		logPrintf(">>>> Ord2Card.4 %s", source)
		return source

	case Portuguese:
		logPrintf(">>>> Ord2Card.0 [word] %s", word)
		if len(word) < 1 {
			return ""
		}
		ordinal, ok := lg.PtOrdinals[word[:len(word)-1]]
		if !ok {
			return ""
		}
		return ordinal
	case English:
		logPrintf(">>>> Ord2Card.1 %s", word)
		plurSuff := strings.HasSuffix(word, "ths")
		singSuff := strings.HasSuffix(word, "th")
		source := ""
		if !(plurSuff || singSuff) {
			if strings.HasSuffix(word, "first") {
				source = strings.ReplaceAll(word, "first", "one")
			} else if strings.HasSuffix(word, "second") {
				source = strings.ReplaceAll(word, "second", "two")
			} else if strings.HasSuffix(word, "third") {
				source = strings.ReplaceAll(word, "third", "three")
			} else {
				logPrintf(">>>> Ord2Card.2 %s", word)
				return ""
			}
		} else {
			if plurSuff {
				source = word[:len(word)-3]
			} else {
				source = word[:len(word)-2]
			}
		}

		if containsKey(lg.RadMap, source) {
			source = lg.RadMap[source]
		} else if strings.HasSuffix(source, "ie") {
			source = source[:len(source)-2] + "y"
		} else if strings.HasSuffix(source, "fif") {
			source = source[:len(source)-1] + "ve"
		} else if strings.HasSuffix(source, "eigh") {
			source = source + "t"
		} else if strings.HasSuffix(source, "nin") {
			source = source + "e"
		}

		if !containsKey(lg.Numbers, source) {
			logPrintf(">>>> Ord2Card.3 %s", source)
			return ""
		}

		logPrintf(">>>> Ord2Card.4 %s", source)
		return source

	default:
		return ""
	}
}

func (lg *Language) NumOrd(digits string, originalWord string) string {
	switch lg.LangCode {
	case French:
		logPrintf(">>>> NumOrd.0 %s", originalWord)

		if containsKey(lg.IrrOrd, originalWord) {
			return lg.IrrOrd[originalWord].One
		}

		if strings.HasSuffix(originalWord, "e") {
			return fmt.Sprintf("%s%s", digits, "ème")
		}
		return fmt.Sprintf("%s%s", digits, "èmes")

	case English:
		logPrintf(">>>> NumOrd.1 %s", originalWord)
		sf := ""
		if strings.HasSuffix(originalWord, "s") {
			sf = originalWord[len(originalWord)-3:]
		} else {
			sf = originalWord[len(originalWord)-2:]
		}

		return fmt.Sprintf("%s%s", digits, sf)

	case Portuguese, Spanish:
		logPrintf(">>>> NumOrd.2 %s", originalWord)
		if strings.HasSuffix(originalWord, "o") {
			return fmt.Sprintf("%sº", digits)
		}

		return fmt.Sprintf("%sª", digits)
	}

	logPrintf(">>>> NumOrd.2 ❌ %s", originalWord)
	return "ERROR"
}

func (lg *Language) Normalize(word string) string {
	switch lg.LangCode {
	case French:
		return strings.ReplaceAll(word, "vingts", "vingt")
	default:
		return word
	}
}

func (lg *Language) NotNumericWord(word string) bool {
	return word == "" || word != lg.DecimalSep && !containsKey(lg.Numbers, word) && !contains(lg.Zero, word)
}

const UsePTOrdinalsMerger = true

var (
	WORDSEP = regexp.MustCompile(`\s*[\.,;\(\)…\[\]:!\?]+\s*|\n`)
	omg     = OrdinalsMerger{}
)

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
		segmentAndPuncts = append(segmentAndPuncts, segmentAndPunct{segment, punct[i]})
	}

	outSegments := []string{}
	for _, sp := range segmentAndPuncts {
		tokens := strings.Split(sp.segment, " ")
		logPrintf("tokens %v", tokens)

		numBuilder := NewWordToDigitParser(lg, relaxed, signed, ordinalThreshold, "")
		lastWord := ""
		inNumber := false
		outTokens := []string{}
		for _, couple := range lookAhead(tokens) {

			logPrintf("✅ [word] %s [ahead] %s", couple.Word, couple.Ahead)

			pushed := numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead))
			if pushed {
				logPrintf("> condition 1: word %s ahead %s", couple.Word, couple.Ahead)
				inNumber = true
			} else if inNumber {
				logPrintf("> condition 2: word %s ahead %s", couple.Word, couple.Ahead)
				outTokens = append(outTokens, numBuilder.GetValue())
				numBuilder = NewWordToDigitParser(lg, relaxed, signed, ordinalThreshold, lastWord)
				inNumber = numBuilder.push(strings.ToLower(couple.Word), strings.ToLower(couple.Ahead))
			}

			if !inNumber {
				logPrintf("> condition 3: word %s ahead %s", couple.Word, couple.Ahead)
				outTokens = append(outTokens, couple.Word)
			}

			lastWord = strings.ToLower(couple.Word)

			logPrintf("... lastWord %s, inNumber %t, outTokens %v", lastWord, inNumber, outTokens)

		}

		logPrintf("---")
		numBuilder.close()
		if numBuilder.GetValue() != "" {
			outTokens = append(outTokens, numBuilder.GetValue())
		}

		outSegments = append(outSegments, strings.Join(outTokens, " "))
		outSegments = append(outSegments, sp.punct)

	}
	text = strings.Join(outSegments, "")

	logPrintf(">>> [text] %s", text)

	// Post-Processing
	if lg.LangCode == Portuguese && UsePTOrdinalsMerger {
		text = omg.MergeCompoundOrdinalsPT(text)
	}

	return text
}
