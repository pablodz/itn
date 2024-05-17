package itn

import (
	"fmt"
	"maps"
)

type LanguageCode int

const (
	Spanish LanguageCode = iota
	English
	French
	Portuguese
)

func NewLanguage(LangCode LanguageCode) (*Language, error) {
	switch LangCode {
	case Spanish:
		l := &Language{
			LangCode: LangCode,
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
				"dieciséis":    16, // with typo
				"veintiséis":   26, // with typo
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
				// "uno", // Telephony first
				"una",
			},
			Relaxed:    map[string]RelaxTuple{},
			Composites: map[string]int{},
		}

		// deep copy from l.multipliers
		l.Numbers = maps.Clone(l.Multipliers)
		maps.Copy(l.Numbers, l.Units)
		maps.Copy(l.Numbers, l.STens)
		maps.Copy(l.Numbers, l.MTens)
		maps.Copy(l.Numbers, l.Hundred)
		maps.Copy(l.Numbers, l.MTens)

		return l, nil

	case English:

		l := &Language{
			LangCode: LangCode,
			Multipliers: map[string]int{
				"hundred":   100,
				"hundreds":  100,
				"thousand":  1000,
				"thousands": 1000,
				"million":   1000000,
				"millions":  1000000,
				"billion":   1000000000,
				"billions":  1000000000,
				"trillion":  1000000000000,
				"trillions": 1000000000000,
			},
			Units: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
				"four":  4,
				"five":  5,
				"six":   6,
				"seven": 7,
				"eight": 8,
				"nine":  9,
			},
			STens: map[string]int{
				"ten":       10,
				"eleven":    11,
				"twelve":    12,
				"thirteen":  13,
				"fourteen":  14,
				"fifteen":   15,
				"sixteen":   16,
				"seventeen": 17,
				"eighteen":  18,
				"nineteen":  19,
			},
			MTens: map[string]int{
				"twenty":  20,
				"thirty":  30,
				"forty":   40,
				"fifty":   50,
				"sixty":   60,
				"seventy": 70,
				"eighty":  80,
				"ninety":  90,
			},
			MTensWSTens: []string{},
			Hundred: map[string]int{
				"hundred":  100,
				"hundreds": 100,
			},
			Sign: map[string]string{
				"plus":  "+",
				"minus": "-",
			},
			Zero: []string{
				"zero",
				"o",
			},
			DecimalSep: "point",
			DecimalSYM: ".",
			AndNums:    []string{},

			And: "and",
			NeverIfAlone: []string{
				"one",
				"o",
			},
			Relaxed: map[string]RelaxTuple{},
			RadMap: map[string]string{
				"fif":   "five",
				"eigh":  "eight",
				"nin":   "nine",
				"twelf": "twelve",
			},
			Composites: map[string]int{},
		}

		for k1, v1 := range l.MTens {
			for k2, v2 := range l.Units {
				l.Composites[fmt.Sprintf("%s-%s", k1, k2)] = v1 + v2
			}
		}

		l.Numbers = maps.Clone(l.Multipliers)
		maps.Copy(l.Numbers, l.Units)
		maps.Copy(l.Numbers, l.STens)
		maps.Copy(l.Numbers, l.MTens)
		maps.Copy(l.Numbers, l.Hundred)
		maps.Copy(l.Numbers, l.Composites)

		return l, nil

	default:
		return nil, fmt.Errorf("Language not implemented")
	}
}
