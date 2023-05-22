package utils

func KMP(pattern []byte, toMatch []byte, wildcard []bool) bool {
	foundMatch := false

	patternLen := len(pattern)

	// compute border
	borders := make([]int, patternLen)
	j, i := 0, 1

	for i < patternLen {
		if pattern[j] == pattern[i] {
			// j+1 chars match
			borders[i] = j + 1
			i++
			j++
		} else if j > 0 { // j follows matching prefix
			j = borders[j-1]
		} else { // no match
			borders[i] = 0
			i++
		}
	}

	// begin matching
	i, j = 0, 0

	for i < len(toMatch) && !foundMatch {
		if pattern[j] == toMatch[i] || wildcard[j] {
			if j == patternLen-1 {
				foundMatch = true
				// result = append(result, i-patternLen+1) // found match
			}
			i++
			j++
		} else if j > 0 {
			j = borders[j-1]
		} else {
			i++
		}
	}

	return foundMatch
}
