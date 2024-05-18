package itn

import "log"

func contains(slice []string, word string) bool {
	for _, v := range slice {
		if v == word {
			return true
		}
	}
	return false
}

func containsKey[T int | string | RelaxTuple](dict map[string]T, key string) bool {
	_, ok := dict[key]
	return ok
}

var DEBUG = false

func SetDebug(debug bool) {
	DEBUG = debug
}

func logPrintf(format string, args ...interface{}) {
	if DEBUG {
		log.Printf(format, args...)
	}
}

func sumInts(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}
