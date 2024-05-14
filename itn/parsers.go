package itn

import (
	"fmt"
	"log"
	"strings"
)

type WordStreamValueParser struct {
	Skip     string
	n000Val  int
	grpVal   int
	lastWord string
	lang     Language
	relaxed  bool
}

func NewWordStreamValueParser(lang Language, relaxed bool) WordStreamValueParser {
	return WordStreamValueParser{
		lang:    lang,
		relaxed: relaxed,
	}
}

func (w *WordStreamValueParser) GetValue() int {
	return w.n000Val + w.grpVal
}

func (w *WordStreamValueParser) groupExpects(word string, update bool) bool {
	expected := false
	if w.lastWord == "" {
		expected = true
	} else if containsKey(w.lang.Units, w.lastWord) && w.grpVal < 10 || containsKey(w.lang.STens, w.lastWord) && w.grpVal < 20 {
		expected = containsKey(w.lang.Hundred, word)
	} else if containsKey(w.lang.MHundreds, w.lastWord) {
		expected = true
	} else if containsKey(w.lang.MTens, w.lastWord) {
		expected = containsKey(w.lang.Units, word) || containsKey(w.lang.STens, word) && contains(w.lang.MTensWSTens, w.lastWord)
	} else if containsKey(w.lang.Hundred, w.lastWord) {
		expected = !containsKey(w.lang.Hundred, word)
	}

	if update {
		w.lastWord = word
	}

	return expected
}

func (w *WordStreamValueParser) isCoefAppliable(coef int) bool {
	if w.lang.Simplify_check_coef_appliable {
		return coef != w.GetValue()
	}

	if coef > w.GetValue() && (w.GetValue() > 0 || coef >= 100) {
		return true
	}

	if coef*1000 <= w.n000Val || coef == 100 && 100 > w.grpVal {
		return (w.grpVal > 0 || coef == 1000 || coef == 100)
	}

	return false
}

func (w *WordStreamValueParser) push(word string, lookAhead string) bool {

	log.Printf("- WordStreamValueParser.push.word %s [ahead] %s", word, lookAhead)

	if word == "" {
		log.Printf(">> WordStreamValueParser.push.condition 0: [word]%s [ahead] %s", word, lookAhead)
		return false
	}

	if word == w.lang.And && contains(w.lang.AndNums, lookAhead) {
		log.Printf(">> WordStreamValueParser.push.condition 1: [word]%s [ahead] %s", word, lookAhead)
		return true
	}

	word = w.lang.Normalize(word)
	if !containsKey(w.lang.Numbers, word) {
		log.Printf(">> WordStreamValueParser.push.condition 2: [word]%s [ahead] %s", word, lookAhead)
		return false
	}

	RELAXED := w.lang.Relaxed
	if containsKey(w.lang.Multipliers, word) {
		log.Printf(">> WordStreamValueParser.push.condition 3: [word]%s [ahead] %s", word, lookAhead)
		coef := w.lang.Multipliers[word]
		log.Printf(">>> WordStreamValueParser.push.coef %d", coef)
		if !w.isCoefAppliable(coef) {
			log.Printf(">> WordStreamValueParser.push.condition 3.1: [word]%s [ahead] %s", word, lookAhead)
			return false
		}

		if coef < 1000 {
			if w.grpVal == 0 {
				w.grpVal = 1
			}
			w.grpVal = w.grpVal * coef
			w.lastWord = ""
			log.Printf(">> WordStreamValueParser.push.condition 3.2: [word]%s [ahead] %s", word, lookAhead)
			return true
		}
		if coef < w.n000Val {
			if w.grpVal == 0 {
				w.grpVal = 1
			}
			w.n000Val = w.n000Val + coef*(w.grpVal)
		} else {
			if w.grpVal == 0 {
				w.grpVal = 1
			}
			w.n000Val = w.GetValue() * coef
		}
		w.grpVal = 0
		w.lastWord = ""
	} else if w.relaxed && containsKey(RELAXED, word) && lookAhead != "" && strings.HasPrefix(RELAXED[word].Zero, lookAhead) && w.groupExpects(RELAXED[word].One, false) {
		log.Printf(">> WordStreamValueParser.push.condition 4: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = RELAXED[word].Zero
		w.grpVal = w.grpVal + w.lang.Numbers[RELAXED[word].One]
	} else if w.Skip != "" && strings.HasPrefix(w.Skip, word) {
		log.Printf(">> WordStreamValueParser.push.condition 5: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = ""
	} else if w.groupExpects(word, true) {
		log.Printf(">> WordStreamValueParser.push.condition 6: [word]%s [ahead] %s", word, lookAhead)
		if containsKey(w.lang.Hundred, word) {
			log.Printf(">> WordStreamValueParser.push.condition 6.1: [word]%s [ahead] %s", word, lookAhead)
			if w.grpVal != 0 {
				w.grpVal = 100 * w.grpVal
			} else {
				w.grpVal = w.lang.Hundred[word]
			}
		} else if containsKey(w.lang.MHundreds, word) {
			log.Printf(">> WordStreamValueParser.push.condition 6.2: [word]%s [ahead] %s", word, lookAhead)
			w.grpVal = w.lang.MHundreds[word]
		} else {
			log.Printf(">> WordStreamValueParser.push.condition 6.3: [word]%s [ahead] %s", word, lookAhead)
			w.grpVal = w.grpVal + w.lang.Numbers[word]
			log.Printf(">>> WordStreamValueParser.push.grpVal %d", w.grpVal)
		}
	} else {
		log.Printf(">> WordStreamValueParser.push.condition 7: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = ""
		return false
	}

	log.Printf(">> WordStreamValueParser.push.condition 8: [word]%s [ahead] %s", word, lookAhead)
	return true
}

type WordToDigitParser struct {
	Lang             *Language
	value            []string
	IntBuilder       WordStreamValueParser
	FracBuilder      WordStreamValueParser
	Signed           bool
	InFrac           bool
	Closed           bool
	Open             bool
	LastWord         string
	OrdinalThreshold int
}

func NewWordToDigitParser(lang *Language, relaxed bool, signed bool, ordinalThreshold int, precedingWord string) *WordToDigitParser {
	return &WordToDigitParser{
		Lang:             lang,
		value:            []string{},
		IntBuilder:       NewWordStreamValueParser(*lang, relaxed),
		FracBuilder:      NewWordStreamValueParser(*lang, relaxed),
		Signed:           signed,
		InFrac:           false,
		Closed:           false,
		Open:             false,
		LastWord:         precedingWord,
		OrdinalThreshold: ordinalThreshold,
	}
}

func (w *WordToDigitParser) GetValue() string {
	return strings.Join(w.value, "")
}

func (w *WordToDigitParser) close() {
	if !w.Closed {
		if w.InFrac && w.FracBuilder.GetValue() > 0 {
			w.value = append(w.value, fmt.Sprint(w.FracBuilder.GetValue()))
		} else if !w.InFrac && w.IntBuilder.GetValue() > 0 {
			w.value = append(w.value, fmt.Sprint(w.IntBuilder.GetValue()))
		}
		w.Closed = true
	}
}

func (w *WordToDigitParser) atStartOfSeq() bool {
	return w.InFrac && w.FracBuilder.GetValue() == 0 || !w.InFrac && w.IntBuilder.GetValue() == 0
}

func (w *WordToDigitParser) atStart() bool {
	return !w.Open
}

func (w *WordToDigitParser) the_push(word string, lookAhead string) bool {
	builder := WordStreamValueParser{}
	log.Printf(">> inFrac %v word %s lookAhead %s", w.InFrac, word, lookAhead)
	if w.InFrac {
		builder = w.FracBuilder
	} else {
		builder = w.IntBuilder
	}
	return builder.push(word, lookAhead)
}

func (w *WordToDigitParser) isAlone(word string, nextWord string) bool {
	return !w.Open && contains(w.Lang.NeverIfAlone, word) && w.Lang.NotNumericWord(nextWord) && w.Lang.NotNumericWord(w.LastWord) && !(nextWord == "" && w.LastWord == "")
}

func (w *WordToDigitParser) push(word string, lookAhead string) bool {

	if w.Closed || w.isAlone(word, lookAhead) {
		log.Printf(">> WordToDigitParser.push.condition 0:[word]%s [ahead] %s", word, lookAhead)
		w.LastWord = word
		return false
	}

	if w.Signed && containsKey(w.Lang.Sign, word) && containsKey(w.Lang.Numbers, lookAhead) && w.atStart() {
		log.Printf(">> WordToDigitParser.push.condition 1:[word]%s [ahead] %s", word, lookAhead)
		w.value = append(w.value, w.Lang.Sign[word])
	} else if contains(w.Lang.Zero, word) && w.atStartOfSeq() && lookAhead != "" && strings.Contains(w.Lang.DecimalSep, lookAhead) {
		log.Printf(">> WordToDigitParser.push.condition 2:[word]%s [ahead] %s", word, lookAhead)
	} else if contains(w.Lang.Zero, word) && w.atStartOfSeq() {
		log.Printf(">> WordToDigitParser.push.condition 3:[word]%s [ahead] %s", word, lookAhead)
		w.value = append(w.value, "0")
	} else if w.the_push(w.Lang.Ord2Card(word), lookAhead) {
		log.Printf(">> WordToDigitParser.push.condition 4:[word]%s [ahead] %s", word, lookAhead)
		value2Add := word
		if w.IntBuilder.GetValue() > w.OrdinalThreshold {
			digits := w.IntBuilder.GetValue()
			if w.InFrac {
				digits = w.FracBuilder.GetValue()
			}
			value2Add = w.Lang.NumOrd(fmt.Sprint(digits), word)
		}
		w.value = append(w.value, value2Add)
		w.Closed = true
	} else if word == w.Lang.DecimalSep || contains(strings.Split(w.Lang.DecimalSep, ","), word) && (containsKey(w.Lang.Numbers, lookAhead) || contains(w.Lang.Zero, lookAhead)) && !w.InFrac {
		log.Printf(">> WordToDigitParser.push.condition 5:[word]%s [ahead] %s", word, lookAhead)
		if w.GetValue() == "" {
			w.value = append(w.value, fmt.Sprint(w.IntBuilder.GetValue()))
		}
		w.value = append(w.value, w.Lang.DecimalSYM)
		w.InFrac = true
	} else if !w.the_push(word, lookAhead) {
		log.Printf(">> WordToDigitParser.push.condition 6:[word] %s [ahead] %s", word, lookAhead)
		if w.Open {
			w.close()
		}
		w.LastWord = word
		return false
	}

	log.Printf(">> WordToDigitParser.push.condition 7:[word] %s [ahead] %s", word, lookAhead)

	w.Open = true
	w.LastWord = word
	return true
}
