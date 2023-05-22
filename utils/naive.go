package utils

func BruteForce(pattern []byte, toMatch []byte, wildcard []bool) bool {
	foundMatch := false
	patternLen := len(pattern)

	i := 0
	maxIdx := len(toMatch) - len(pattern) + 1

	for i < maxIdx && !foundMatch {
		equal := true
		j := 0

		for j < patternLen && equal {
			equal = pattern[j] == toMatch[i+j] || wildcard[j]
			j++
		}

		if equal {
			foundMatch = true
		}
		i++
	}

	return foundMatch
}
