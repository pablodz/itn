package itn

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var SegmentBreak = regexp.MustCompile(`\s*[\.,;\(\)…\[\]:!\?]+\s*`)

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
		segmentAndPuncts = append(segmentAndPuncts, segmentAndPunct{segment, punct[i]})
	}

	outSegments := []string{}
	ordinal := 0
	gender := ""

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

		logPrintf("> MergeCompoundOrdinalsPT.1.tokens %v [tokens2] %v", tokens, tokens2)

		for pointer < len(tokens) {
			token := tokens[pointer]

			if om.isOrdinal(token) {
				logPrintf("> MergeCompoundOrdinalsPT.2.1.token <%s> [tokens2] %v", token, tokens2)
				currentIsOrdinal = true
				seq = append(seq, om.getCardinal(token))
				gender = om.getGender(token)
			} else {
				if !currentIsOrdinal {
					logPrintf("> MergeCompoundOrdinalsPT.4.token %s [tokens2] %v", token, tokens2)
					tokens2 = append(tokens2, token)
				} else {
					logPrintf("> MergeCompoundOrdinalsPT.5.token %s [tokens2] %v", token, tokens2)
					logPrintf("> MergeCompoundOrdinalsPT.5.seq %v", seq)
					ordinal = sumInts(seq)
					logPrintf("> MergeCompoundOrdinalsPT.5.ordinal %d", ordinal)
					tokens2 = append(tokens2, fmt.Sprintf("%s%s", strconv.Itoa(ordinal), gender))
					tokens2 = append(tokens2, token)
					seq = []int{}
					currentIsOrdinal = false
					logPrintf("> MergeCompoundOrdinalsPT.5.1.token %s [tokens2] %v", token, tokens2)

				}
			}

			pointer++

		}

		if currentIsOrdinal {
			logPrintf("> MergeCompoundOrdinalsPT.6.seq %v [tokens2] %v", seq, tokens2)
			ordinal = sumInts(seq)
			tokens2 = append(tokens2, fmt.Sprintf("%s%s", strconv.Itoa(ordinal), gender))
		}

		tokens2 = om.text2NumStyle(tokens2)
		sp.segment = strings.Join(tokens2, " ") + sp.punct
		outSegments = append(outSegments, sp.segment)

	}
	text = strings.Join(outSegments, "")
	logPrintf("> MergeCompoundOrdinalsPT.7.text %s", text)

	return text
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
	runes := []rune(token)
	numPart := string(runes[:len(runes)-1]) // Extract the part of the string before the last character
	logPrintf(">>>> [getCardinal] token[:-1] %s", numPart)
	out, err := strconv.Atoi(numPart)
	if err != nil {
		switch numPart {
		case "primeir":
			out = 1
		case "segund":
			out = 2
		case "terceir":
			out = 3
		default:
			out = 0
		}
	}
	logPrintf(">>>> [getCardinal] %s -> %d", token, out)
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

type SubRegex struct {
	Pattern     *regexp.Regexp
	Replacement string
}

var subRegexes = []SubRegex{
	{regexp.MustCompile(`1\s`), "um "},
	{regexp.MustCompile(`2\s`), "dois"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])1[º°]([^a-zA-Z0-9]|$)`), "primeiro"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])2[º°]([^a-zA-Z0-9]|$)`), "segundo"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])3[º°]([^a-zA-Z0-9]|$)`), "terceiro"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])1ª([^a-zA-Z0-9]|$)`), "primeira"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])2ª([^a-zA-Z0-9]|$)`), "segunda"},
	{regexp.MustCompile(`(^|[^a-zA-Z0-9])3ª([^a-zA-Z0-9]|$)`), "terceira"},
}

func (om *OrdinalsMerger) text2NumStyle(tokens []string) []string {
	logPrintf(">>>>>>> [Tokens] --- %v", tokens)
	for i, token := range tokens {
		for _, sr := range subRegexes {
			v := sr.Pattern.ReplaceAllString(token, sr.Replacement)
			if token != v {
				logPrintf(">>>>>>>>>>>>>>>>>>>>>>>>> [Tokens] !!![%d] %s -> %s", i, tokens[i], v)
			}
			token = v
		}
		tokens[i] = token
	}
	logPrintf(">>>>>>> [Tokens] --- %v", tokens)
	return tokens
}
