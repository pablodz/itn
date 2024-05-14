package itn

import (
	"fmt"
	"strings"
)

type SpanishLanguage struct {
	*Language
	Composites map[string]int
}

func NewSpanishLanguage() *SpanishLanguage {
	l := &SpanishLanguage{}
	l.Language = &Language{}

	l.Multipliers = map[string]int{
		"mil":      1000,
		"miles":    1000,
		"millon":   1000000,
		"millón":   1000000,
		"millones": 1000000,
	}

	l.Units = map[string]int{
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
	}

	l.STens = map[string]int{
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
	}

	l.MTens = map[string]int{
		"treinta":   30,
		"cuarenta":  40,
		"cincuenta": 50,
		"sesenta":   60,
		"setenta":   70,
		"ochenta":   80,
		"noventa":   90,
	}

	l.MTensWSTens = []string{}

	l.Hundred = map[string]int{
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
	}

	l.Composites = map[string]int{}

	// deep copy from l.multipliers
	l.Numbers = l.Multipliers
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
	for k, v := range l.Composites {
		l.Numbers[k] = v
	}

	l.Sign = map[string]string{
		"mas":   "+",
		"menos": "-",
	}

	l.Zero = []string{
		"cero",
	}

	l.DecimalSep = "coma"
	l.DecimalSYM = "."

	l.AndNums = []string{
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
	}

	l.And = "y"

	l.NeverIfAlone = []string{
		"un",
		"uno",
		"una",
	}

	return l
}

func (lg *SpanishLanguage) Ord2Card(word string) string {
	return word
}

func (lg *SpanishLanguage) NumOrd(digits string, originalWord string) string {
	if strings.HasSuffix(originalWord, "o") {
		return fmt.Sprintf("%sº", digits)
	}
	return fmt.Sprintf("%sª", digits)
}
