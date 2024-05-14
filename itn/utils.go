package itn

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
