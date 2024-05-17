package itn

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var SegmentBreak = regexp.MustCompile(`\s*[\.,;\(\)…\[\]:!\?]+\s*`)

type SubRegex struct {
	regex       *regexp.Regexp
	replacement string
}

// Initialize the regexes and replacements
var subRegexes = []SubRegex{
	{regexp.MustCompile(`1\s`), "um "},
	{regexp.MustCompile(`2\s`), "dois"},
	{regexp.MustCompile(`\b1[º°]\b`), "primeiro"},
	{regexp.MustCompile(`\b2[º°]\b`), "segundo"},
	{regexp.MustCompile(`\b3[º°]\b`), "terceiro"},
	{regexp.MustCompile(`\b1ª\b`), "primeira"},
	{regexp.MustCompile(`\b2ª\b`), "segunda"},
	{regexp.MustCompile(`\b3ª\b`), "terceira"},
}

type OrdinalsMerger struct{}

func NewOrdinalsMerger() *OrdinalsMerger {
	return &OrdinalsMerger{}
}

func (om *OrdinalsMerger) MergeCompoundOrdinalsPT(text string) string {
	segments := SegmentBreak.Split(text, -1)
	punct := SegmentBreak.FindAllString(text, -1)
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
		tokens := []string{}
		for _, t := range strings.Split(sp.segment, " ") {
			if len(t) > 0 {
				tokens = append(tokens, t)
			}
		}

		pointer := 0
		tokens2 := []string{}
		currentIsOrdinal := false
		seq := []int{}
		gender := ""
		ordinal := 0

		for pointer < len(tokens) {
			token := tokens[pointer]
			if om.isOrdinal(token) {
				currentIsOrdinal = true
				seq = append(seq, om.getCardinal(token))
				gender = om.getGender(token)
			} else {
				if !currentIsOrdinal {
					tokens2 = append(tokens2, token)
				} else {
					for _, s := range seq {
						ordinal = ordinal + s
					}
					tokens2 = append(tokens2, fmt.Sprintf("%s%s", strconv.Itoa(ordinal), gender))
					tokens2 = append(tokens2, token)
					seq = []int{}
					currentIsOrdinal = false
				}
			}

			pointer++

		}

		if currentIsOrdinal {
			for _, s := range seq {
				ordinal = ordinal + s
			}
			tokens2 = append(tokens2, fmt.Sprintf("%s%s", strconv.Itoa(ordinal), gender))
		}

		tokens2 = om.text2NumStyle(tokens2)
		sp.segment = strings.Join(tokens2, " ") + sp.punct
		outSegments = append(outSegments, sp.segment)

	}

	return strings.Join(outSegments, "")
}

func (om *OrdinalsMerger) isOrdinal(token string) bool {
	out := false
	if len(token) > 1 && (strings.Contains(token, "º") || strings.Contains(token, "°") || strings.Contains(token, "ª")) {
		out = true
	}
	if contains([]string{
		"primeiro",
		"primeira",
		"segundo",
		"segunda",
		"terceiro",
		"terceira",
	}, token) {
		out = true
	}

	return out
}

func (om *OrdinalsMerger) getCardinal(token string) int {
	out := 0
	token = strings.TrimSpace(token)
	if len(token) < 2 {
		return out
	}
	numPart := token[:len(token)-1]
	out, err := strconv.Atoi(numPart)
	if err != nil {
		switch numPart {
		case "primeir":
			out = 1
		case "segund":
			out = 2
		case "terceir":
			out = 3
		}
	}
	return out
}

func (om *OrdinalsMerger) getGender(token string) string {
	gender := string(token[len(token)-1])
	if gender == "a" {
		gender = "ª"
	}
	if gender == "o" {
		gender = "º"
	}
	return gender
}

func (om *OrdinalsMerger) text2NumStyle(tokens []string) []string {
	for i, token := range tokens {
		for _, r := range subRegexes {
			token = r.regex.ReplaceAllString(token, r.replacement)
		}
		tokens[i] = token
	}
	return tokens
}
