package main

import (
	"fmt"
//	"math"
	"sort"
	"strings"
)


//
// NGramSet
//
type NGramSet struct {
	N int
	Total int
	NGrams map[string]int // pattern => count
}

func (set *NGramSet) Add(ngram string) {
	set.NGrams[ngram] += 1
	set.Total += 1
}

func (set *NGramSet) Merge(otherSet NGramSet) {
	for key := range otherSet.NGrams {
		set.NGrams[key] += otherSet.NGrams[key]
		set.Total += otherSet.NGrams[key]
	}
}

func (set NGramSet) Print() {
	orderedKeys := set.OrderedNgrams()
	for _, key := range orderedKeys {
		fmt.Printf("%s --> %d [%2.2f%%]\n", key, set.NGrams[key], set.GetFreq(key) * 100.0)
	}

}

func (set NGramSet) GetFreq(ngram string) float64 {
	return float64(set.NGrams[ngram]) / float64(set.Total)
}

func (set *NGramSet) OrderedNgrams() []string {
	counts := make([]int, 0, len(set.NGrams))
	for _, count := range set.NGrams {
		counts = append(counts, count)
	}
	sort.Ints(counts)

	ngramsByCount := make(map[int][]string, 0)
	for ngram, count := range set.NGrams {
		if ngramsByCount[count] == nil {
			ngramsByCount[count] = make([]string, 0)
		}
		ngramsByCount[count] = append(ngramsByCount[count], ngram)
	}


	results := make([]string, 0)

	includedCounts := make(map[int]bool, 0)

	for _, count := range counts {
		if includedCounts[count] {
			continue
		}
		includedCounts[count] = true

		results = append(results, ngramsByCount[count]...)
	}

	return results
}

func CreateNGramSet(n int) *NGramSet {
	return &NGramSet{n, 0, make(map[string]int)}
}



//
// NGramUniverse
//
type NGramUniverse struct {
	MaxN int
	NgramsByN map[int]*NGramSet
	Alphabet map[string]bool
}

func (universe NGramUniverse) AddString(s string) {
	n := 1
	for n <= universe.MaxN {
		newNGramSet := CountNGrams(s, n)
		universe.NgramsByN[n].Merge(newNGramSet)
		n++
	}

	for _, c := range s {
		universe.Alphabet[string(c)] = true
	}
}

func (universe NGramUniverse) GetAlphabetCharacters() []string {
	chars := make([]string, 0, 0)

	for char := range universe.Alphabet {
		chars = append(chars, char)
	}

	return chars
}

func (universe NGramUniverse) Generate(m int) string {
	return ""
}

func (universe NGramUniverse) Print() {
	n := 1
	for n <= universe.MaxN {
		fmt.Println()
		universe.NgramsByN[n].Print()
		n++
	}
}

func (universe *NGramUniverse) PrintProbabilitiesOfCompletion(s string)  {
	alphabetCharacters := universe.GetAlphabetCharacters()

	fmt.Printf("\n\n%s\n-----------\n", s)

	for _, char := range alphabetCharacters {
		fmt.Printf("%s -> %s   %2.2f%%\n", s, char, 100.0 * universe.GetProbabilityOfCompletion(s, char))
	}
}

func (universe *NGramUniverse) ProbableCompletion(s string) (next string, prob float64)  {
	alphabetCharacters := universe.GetAlphabetCharacters()
	bestCompletion := ""
	bestCompletionProb := 0.0

	for _, char := range alphabetCharacters {
		prob := universe.GetProbabilityOfCompletion(s, char)
		if prob > bestCompletionProb {
			bestCompletion = char
			bestCompletionProb = prob
		}
	}

	return bestCompletion, bestCompletionProb
}

func (universe *NGramUniverse) GetProbabilityOfCompletion(s string, candidate string) float64 {
	if len(s) >= universe.MaxN {
		panic("string too long to make prediction")
	}

	total := 0

	set := universe.NgramsByN[len(s)]
	for ngram, count := range set.NGrams {
		if strings.HasPrefix(ngram, s) {
			total += count
		}
	}

	if total == 0 {
		fmt.Println("[WARN] no ngrams have prefix " + s)
		return 0
	}

	return float64(universe.NgramsByN[len(s) + 1].NGrams[s + candidate])/float64(total)
}

func (universe NGramUniverse) GetCompletionScore(s string, nextChar string) float64 {
	completion := s + nextChar
	probs := make([]float64, 0)
	n := 1
	for n <= universe.MaxN && n <= len(completion) {
		end := len(s) - n
		if end < 0 {
			end = 0
		}

		ngramStart := len(completion) - n
		ngram := completion[ngramStart:]
		prob := universe.NgramsByN[n].GetFreq(ngram)

		probs = append(probs, prob)
		n++
	}

	// TODO: consider alternate ways to combine these averages
	// For example, more weight to the longer ngrams?
	score := 0.0
	for _, p := range probs {
		score += p // * math.Pow(float64(i + 1), 1)
	}

	return score / float64(len(probs))
}

func (universe *NGramUniverse) GenerateStringOfLength(m int) string {
	s := ""

	for len(s) < m {
		nextChar := universe.GenerateNextCharacter(s)
		s = s + nextChar
	}

	return s
}

func (universe *NGramUniverse) GenerateNextCharacter(s string) string {
	alphabetCharacters := universe.GetAlphabetCharacters()

	bestCompletion := alphabetCharacters[0]
	bestProb := 0.0

	for _, char := range alphabetCharacters {
		prob := universe.GetCompletionScore(s, char)
		if prob > bestProb {
			bestCompletion = char
			bestProb = prob
		}
	}

	return bestCompletion
}




func CreateUniverse(maxN int) NGramUniverse {
	universe := NGramUniverse{}
	universe.MaxN = maxN
	universe.NgramsByN = make(map[int]*NGramSet)
	universe.Alphabet = make(map[string]bool)

	n := 1
	for n <= universe.MaxN {
		universe.NgramsByN[n] = CreateNGramSet(n)
		n++
	}

	return universe
}


//
// Functions
//

func CountNGrams(s string, n int) NGramSet {
	ngrams := ExtractNGrams(s, n)
	set := NGramSet{n, 0, make(map[string]int)}

	for _, ngram := range ngrams {
		set.Add(ngram)
	}

	return set
}

func ExtractNGrams(s string, n int) []string {
	ngrams := make([]string, 0)
	for i, _ := range s {
		start := i
		end := i + n

		if end <= len(s) {
			ngram := s[start:end]
			ngrams = append(ngrams, ngram)
		}
	}

	return ngrams
}
