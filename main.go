package main

import (
	"encoding/hex"
	"fmt"
	"makalah/utils"
	"strconv"
	"strings"
	"time"
)

type matcher func([]byte, []byte, []bool) bool

func testPatterns(source []byte, patterns [][]byte, wildcards [][]bool, tester matcher, name string) []bool {
	result := make([]bool, 0)

	var totalElapsed time.Duration = 0

	fmt.Printf("%s\n", name)

	i := 0
	for i < len(patterns) {
		start := time.Now()
		result = append(result, tester(patterns[i], source, wildcards[i]))
		elapsed := time.Since(start)
		fmt.Printf("Pattern %s took %s\n", hex.EncodeToString(patterns[i]), elapsed)
		i++
		totalElapsed += elapsed
	}

	fmt.Printf("Total elapsed %s\n", totalElapsed)

	// fmt.Print("result ")
	// for _, res := range result {
	// 	fmt.Print(strconv.FormatBool(res))
	// 	fmt.Print(" ")
	// }
	fmt.Print("\n")

	return result
}

func analyzeFile(url string, password string, patterns [][]byte, wildcards [][]bool) [][]bool {
	fileToAnalyze := utils.ReadFile(url, password)
	fmt.Printf("File size %d bytes\n", len(fileToAnalyze))

	bruteResult := testPatterns(fileToAnalyze, patterns, wildcards, utils.BruteForce, "Brute Force")
	bmResult := testPatterns(fileToAnalyze, patterns, wildcards, utils.BM, "Boyer Moore")
	kmpResult := testPatterns(fileToAnalyze, patterns, wildcards, utils.KMP, "KMP")

	result := make([][]bool, 0)
	result = append(result, bruteResult)
	result = append(result, bmResult)
	result = append(result, kmpResult)

	return result
}

func strtoByte(input string) ([]byte, []bool) {
	bytes := make([]byte, 0)
	wildcard := make([]bool, 0)

	splitted := strings.Fields(input)

	for _, each := range splitted {
		if strings.Contains(each, "?") {
			bytes = append(bytes, 0x00)
			wildcard = append(wildcard, true)
		} else {
			ui, err := strconv.ParseUint(each, 16, 8)

			if err != nil {
				panic("Cannot parse")
			}

			bytes = append(bytes, byte(ui))
			wildcard = append(wildcard, false)
		}
	}

	return bytes, wildcard
}

func testFile(name string, hash string, strPatterns []string) {

	fmt.Printf("Testing %s with hash %s\n", name, hash)

	patterns := make([][]byte, 0)
	wildcards := make([][]bool, 0)

	for _, pat := range strPatterns {
		p, w := strtoByte(pat)
		patterns = append(patterns, p)
		wildcards = append(wildcards, w)
	}

	analyzeFile(hash, "infected", patterns, wildcards)
	fmt.Println("")
}

func main() {
	// /*TESTING TRICKBOT
	// YARA PATTERN: https://malpedia.caad.fkie.fraunhofer.de/details/win.trickbot
	// 				sequence 0 and 2
	// Sample source: https://bazaar.abuse.ch/sample/6857d1f39cca9db038166c47cb8742ee8f03601bf4e23394a5b94c35ae3d1726/
	// */
	trickbotPatterns := []string{
		"83 e0 70 83 c0 10 eb 25 a9 00 00 00 40 74 11 25 00 00 00 80",
		"f7 d8 1b c0 83 e0 20 83 c0 20 eb 36 25 00 00 00 80 f7 d8",
	}
	testFile("Trickbot", "6857d1f39cca9db038166c47cb8742ee8f03601bf4e23394a5b94c35ae3d1726", trickbotPatterns)

	// /*TESTING WANNARCRY
	// YARA PATTERN: https://malpedia.caad.fkie.fraunhofer.de/details/win.wannacryptor
	// 				sequence 0, 1, and 2
	// Sample source: https://bazaar.abuse.ch/sample/561d7f05055800d3eb9d9e150969e2c84a71dc82a362fb3e1a224af420e53b35/
	// */
	wannaCryPatterns := []string{
		"83 ec 1c 56 8b f1 8a 46 5a 84 c0 75 6f 8a 46 58",
		"c2 0c 00 8b ce e8 ?? ?? ?? ?? 84 c0",
		"74 13 3c 2f 74 04 3c 5c 75 03 8d 6e 01 8a 46 01",
	}
	testFile("WannaCry", "561d7f05055800d3eb9d9e150969e2c84a71dc82a362fb3e1a224af420e53b35", wannaCryPatterns)

	/*TESTING REMCOS
	YARA PATTERN: https://malpedia.caad.fkie.fraunhofer.de/details/win.remcos
					sequence 0, and 1
	Sample source: https://bazaar.abuse.ch/sample/6c0f5a9bf9bfd84be91f3d84335b63ac95ac2b227fedc5de439971577328ac30/
	*/
	remcosPatterns := []string{
		"ff 15 ?? ?? ?? ?? 85 c0 74 10 6a 00 ff 35 ?? ?? ?? ??",
		"6a 09 ff 35 ?? ?? ?? ?? ff 15 ?? ?? ?? ?? ff 35 ?? ?? ?? ??",
		"74 10 6a 00 ff 35 ?? ?? ?? ?? ff 15 ?? ?? ?? ??",
	}
	testFile("Remcos", "6c0f5a9bf9bfd84be91f3d84335b63ac95ac2b227fedc5de439971577328ac30", remcosPatterns)

	/*TESTING NANOCORE
	YARA PATTERN: https://malpedia.caad.fkie.fraunhofer.de/details/win.nanocore
					key
	Sample source: https://bazaar.abuse.ch/sample/920f8094f26c8540082ba329bbdbdba2ba082d4f6b8427103f6acd2ae66ef89c/
	*/
	nanocorePatterns := []string{
		"43 6f 24 cb 95 30 38 39",
	}
	testFile("Nanocore", "920f8094f26c8540082ba329bbdbdba2ba082d4f6b8427103f6acd2ae66ef89c", nanocorePatterns)
}
