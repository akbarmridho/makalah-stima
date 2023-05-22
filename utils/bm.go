package utils

func BadCharHeuristic(pattern []byte, patLen int) []int {
	ret := make([]int, 0)

	for i := 0; i < 256; i++ { //intialize all chars to be of index -1
		ret = append(ret, -1)
	}
	for i := 0; i < patLen; i++ {
		ret[int(pattern[i])] = i // updates the last index of all the chars in the pattern
	}
	return ret
}

func BM(patternStr []byte, toMatchStr []byte, wildcard []bool) bool {
	// make both the pattern and to match string into lowercase
	var patLen int = len(patternStr)
	var textLen int = len(toMatchStr)
	// ret := make([]int, 0)
	foundMatch := false

	if patLen <= textLen {
		var badChar []int = BadCharHeuristic(patternStr, patLen) // calculate the bad char heuristic
		i := 0
		for i <= textLen-patLen && !foundMatch {
			j := patLen - 1
			// check from backwords
			for j >= 0 && (patternStr[j] == toMatchStr[i+j] || wildcard[j]) {
				j--
			}

			if j < 0 { // if it is a matching string
				// ret = append(ret, i)
				foundMatch = true
				// shift the pattern so the next char is alligned with its last occurence in the pattern
				if i+patLen < textLen {
					i += patLen - badChar[toMatchStr[i+patLen]]
				} else {
					i += 1
				}
			} else {
				// shift the pattern so that the mismatched character aligns with the last occurence of it in pattern
				if diff := j - badChar[toMatchStr[i+j]]; diff > 1 {
					i += diff
				} else {
					i += 1
				}
			}
		}
	}
	return foundMatch
}
