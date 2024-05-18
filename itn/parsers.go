package itn

import (
	"fmt"
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

func NewWordStreamValueParser(lang Language, relaxed bool) *WordStreamValueParser {
	return &WordStreamValueParser{
		lang:    lang,
		relaxed: relaxed,
	}
}

func (w *WordStreamValueParser) GetValue() int {
	logPrintf("+ WordStreamValueParser.GetValue")
	return w.n000Val + w.grpVal
}

func (w *WordStreamValueParser) groupExpects(word string, update bool) bool {
	logPrintf("+ WordStreamValueParser.groupExpects.word %s [lastWord] %s [update] %t", word, w.lastWord, update)
	expected := false
	if w.lastWord == "" {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 0: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		expected = true
	} else if containsKey(w.lang.Units, w.lastWord) && w.grpVal < 10 || containsKey(w.lang.STens, w.lastWord) && w.grpVal < 20 {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 1: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		expected = containsKey(w.lang.Hundred, word)
	} else if containsKey(w.lang.MHundreds, w.lastWord) {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 2: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		expected = true
	} else if containsKey(w.lang.MTens, w.lastWord) {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 3: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		expected = containsKey(w.lang.Units, word) || containsKey(w.lang.STens, word) && contains(w.lang.MTensWSTens, w.lastWord)
	} else if containsKey(w.lang.Hundred, w.lastWord) {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 4: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		expected = !containsKey(w.lang.Hundred, word)
	}

	if update {
		logPrintf(">> WordStreamValueParser.groupExpects.condition 5: [word]%s [lastWord] %s [update] %t", word, w.lastWord, update)
		w.lastWord = word
	}

	return expected
}

func (w *WordStreamValueParser) isCoefAppliable(coef int) bool {
	logPrintf("+ WordStreamValueParser.isCoefAppliable.coef %d", coef)
	if w.lang.Simplify_check_coef_appliable {
		logPrintf(">> WordStreamValueParser.isCoefAppliable.condition 0: [coef] %d", coef)
		return coef != w.GetValue()
	}

	if coef > w.GetValue() && (w.GetValue() > 0 || coef >= 100) {
		logPrintf(">> WordStreamValueParser.isCoefAppliable.condition 1: [coef] %d", coef)
		return true
	}

	if coef*1000 <= w.n000Val || coef == 100 && 100 > w.grpVal {
		logPrintf(">> WordStreamValueParser.isCoefAppliable.condition 2: [coef] %d", coef)
		return w.grpVal > 0 || coef == 1000 || coef == 100
	}

	logPrintf(">> WordStreamValueParser.isCoefAppliable.condition 3: [coef] %d", coef)

	return false
}

func (w *WordStreamValueParser) push(word string, lookAhead string) bool {
	logPrintf("+ WordStreamValueParser.push.word %s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)

	if word == "" {
		logPrintf(">> WordStreamValueParser.push.condition 0: [word]%s [ahead] %s", word, lookAhead)
		return false
	}

	if word == w.lang.And && contains(w.lang.AndNums, lookAhead) {
		logPrintf(">> WordStreamValueParser.push.condition 1: [word]%s [ahead] %s", word, lookAhead)
		return true
	}

	word = w.lang.Normalize(word)
	if !containsKey(w.lang.Numbers, word) {
		logPrintf(">> WordStreamValueParser.push.condition 2: [word]%s [ahead] %s", word, lookAhead)
		return false
	}

	RELAXED := w.lang.Relaxed
	if containsKey(w.lang.Multipliers, word) {
		logPrintf(">> WordStreamValueParser.push.condition 3: [word]%s [ahead] %s", word, lookAhead)
		coef := w.lang.Multipliers[word]
		logPrintf(">>> WordStreamValueParser.push.coef %d", coef)
		if !w.isCoefAppliable(coef) {
			logPrintf(">> WordStreamValueParser.push.condition 3.1: [word]%s [ahead] %s", word, lookAhead)
			return false
		}

		if coef < 1000 {
			value := w.grpVal
			if value == 0 {
				value = 1
			}
			w.grpVal = value * coef
			w.lastWord = ""
			logPrintf(">> WordStreamValueParser.push.condition 3.2: [word]%s [ahead] %s", word, lookAhead)
			return true
		}
		if coef < w.n000Val {
			value := w.grpVal
			if value == 0 {
				value = 1
			}
			w.n000Val = w.n000Val + coef*(value)
		} else {
			value := w.GetValue()
			if value == 0 {
				value = 1
			}
			logPrintf(">>> WordStreamValueParser.push.condition 3.3: [value] %d [coef] %d", value, coef)
			w.n000Val = value * coef
		}
		w.grpVal = 0
		w.lastWord = ""
	} else if w.relaxed && containsKey(RELAXED, word) && lookAhead != "" && strings.HasPrefix(RELAXED[word].Zero, lookAhead) && w.groupExpects(RELAXED[word].One, false) {
		logPrintf(">> WordStreamValueParser.push.condition 4: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = RELAXED[word].Zero
		w.grpVal = w.grpVal + w.lang.Numbers[RELAXED[word].One]
	} else if w.Skip != "" && strings.HasPrefix(w.Skip, word) {
		logPrintf(">> WordStreamValueParser.push.condition 5: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = ""
	} else if w.groupExpects(word, true) {
		logPrintf(">> WordStreamValueParser.push.condition 6: [word]%s [ahead] %s", word, lookAhead)
		if containsKey(w.lang.Hundred, word) {
			logPrintf(">> WordStreamValueParser.push.condition 6.1: [word]%s [ahead] %s", word, lookAhead)
			if w.grpVal != 0 {
				w.grpVal = 100 * w.grpVal
			} else {
				w.grpVal = w.lang.Hundred[word]
			}
		} else if containsKey(w.lang.MHundreds, word) {
			logPrintf(">> WordStreamValueParser.push.condition 6.2: [word]%s [ahead] %s", word, lookAhead)
			w.grpVal = w.lang.MHundreds[word]
		} else {
			logPrintf(">> WordStreamValueParser.push.condition 6.3: [word]%s [ahead] %s", word, lookAhead)
			w.grpVal = w.grpVal + w.lang.Numbers[word]
			logPrintf(">>> WordStreamValueParser.push.grpVal %d", w.grpVal)
		}
	} else {
		logPrintf(">> WordStreamValueParser.push.condition 7: [word]%s [ahead] %s", word, lookAhead)
		w.Skip = ""
		return false
	}

	logPrintf(">> WordStreamValueParser.push.condition 8: [word]%s [ahead] %s", word, lookAhead)
	return true
}

type WordToDigitParser struct {
	Lang             Language
	the_value        []string
	IntBuilder       *WordStreamValueParser
	FracBuilder      *WordStreamValueParser
	Signed           bool
	InFrac           bool
	Closed           bool
	Open             bool
	lastWord         string
	OrdinalThreshold int
}

func NewWordToDigitParser(lang Language, relaxed bool, signed bool, ordinalThreshold int, precedingWord string) WordToDigitParser {
	return WordToDigitParser{
		Lang:             lang,
		the_value:        []string{},
		IntBuilder:       NewWordStreamValueParser(lang, relaxed),
		FracBuilder:      NewWordStreamValueParser(lang, relaxed),
		Signed:           signed,
		InFrac:           false,
		Closed:           false,
		Open:             false,
		lastWord:         precedingWord,
		OrdinalThreshold: ordinalThreshold,
	}
}

func (w *WordToDigitParser) GetValue() string {
	logPrintf("+ WordToDigitParser.GetValue")
	return strings.Join(w.the_value, "")
}

func (w *WordToDigitParser) close() {
	logPrintf("+ WordToDigitParser.close")
	if !w.Closed {
		if w.InFrac && w.FracBuilder.GetValue() != 0 {
			logPrintf(">> WordToDigitParser.close.condition 0: adding FracBuilder %d", w.FracBuilder.GetValue())
			w.the_value = append(w.the_value, fmt.Sprint(w.FracBuilder.GetValue()))
		} else if !w.InFrac && w.IntBuilder.GetValue() != 0 {
			logPrintf(">> WordToDigitParser.close.condition 1: adding IntBuilder %d", w.IntBuilder.GetValue())
			w.the_value = append(w.the_value, fmt.Sprint(w.IntBuilder.GetValue()))
		}
		w.Closed = true
	}
}

func (w *WordToDigitParser) atStartOfSeq() bool {
	logPrintf(">> WordToDigitParser.atStartOfSeq")
	return w.InFrac && w.FracBuilder.GetValue() == 0 || !w.InFrac && w.IntBuilder.GetValue() == 0
}

func (w *WordToDigitParser) atStart() bool {
	logPrintf(">> WordToDigitParser.atStart")
	return !w.Open
}

func (w *WordToDigitParser) the_push(word string, lookAhead string) bool {
	logPrintf("🌀 >> inFrac %v [word] %s [lookAhead] %s [lastWord] %s", w.InFrac, word, lookAhead, w.lastWord)
	if w.InFrac {
		builder := w.FracBuilder
		return builder.push(word, lookAhead)
	} else {
		builder := w.IntBuilder
		return builder.push(word, lookAhead)
	}
}

func (w *WordToDigitParser) isAlone(word string, nextWord string) bool {
	return !w.Open && contains(w.Lang.NeverIfAlone, word) && w.Lang.NotNumericWord(nextWord) && w.Lang.NotNumericWord(w.lastWord) && !(nextWord == "" && w.lastWord == "")
}

func (w *WordToDigitParser) push(word string, lookAhead string) bool {
	if w.Closed || w.isAlone(word, lookAhead) {
		logPrintf(">> WordToDigitParser.push.condition 0:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		w.lastWord = word
		return false
	}

	if w.Signed && containsKey(w.Lang.Sign, word) && containsKey(w.Lang.Numbers, lookAhead) && w.atStart() {
		logPrintf(">> WordToDigitParser.push.condition 1:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		w.the_value = append(w.the_value, w.Lang.Sign[word])
	} else if contains(w.Lang.Zero, word) && w.atStartOfSeq() && lookAhead != "" && strings.Contains(w.Lang.DecimalSep, lookAhead) {
		logPrintf(">> WordToDigitParser.push.condition 2:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
	} else if contains(w.Lang.Zero, word) && w.atStartOfSeq() {
		logPrintf(">> WordToDigitParser.push.condition 3:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		w.the_value = append(w.the_value, "0")
	} else if w.the_push(w.Lang.Ord2Card(word), lookAhead) {
		logPrintf(">> WordToDigitParser.push.condition 4:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		value2Add := ""
		if w.IntBuilder.GetValue() > w.OrdinalThreshold {
			digits := 0
			if w.InFrac {
				digits = w.FracBuilder.GetValue()
			} else {
				digits = w.IntBuilder.GetValue()
			}
			value2Add = w.Lang.NumOrd(fmt.Sprint(digits), word)
		} else {
			value2Add = word
		}
		w.the_value = append(w.the_value, value2Add)
		w.Closed = true
	} else if (word == w.Lang.DecimalSep || contains(strings.Split(w.Lang.DecimalSep, ","), word)) && (containsKey(w.Lang.Numbers, lookAhead) || contains(w.Lang.Zero, lookAhead)) && !w.InFrac {
		logPrintf(">> WordToDigitParser.push.condition 5:[word]%s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		if w.GetValue() == "" {
			w.the_value = append(w.the_value, fmt.Sprint(w.IntBuilder.GetValue()))
		}
		w.the_value = append(w.the_value, w.Lang.DecimalSYM)
		w.InFrac = true
	} else if !w.the_push(word, lookAhead) {
		logPrintf(">> WordToDigitParser.push.condition 6:[word] %s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)
		if w.Open {
			w.close()
		}
		w.lastWord = word
		return false
	}

	logPrintf(">> WordToDigitParser.push.condition 7:[word] %s [ahead] %s [lastWord] %s", word, lookAhead, w.lastWord)

	w.Open = true
	w.lastWord = word
	return true
}
