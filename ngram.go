package main

import (
	"fmt"
//	"math"
)


//
// NGramSet
//
type NGramSet struct {
	N int
	Total int64
	NGrams map[string]int64 // pattern => count
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
	for key := range set.NGrams {
		fmt.Printf("%s --> %d\n", key, set.NGrams[key])
	}
}

func (set NGramSet) GetFreq(ngram string) float64 {
	return float64(set.NGrams[ngram]) / float64(set.Total)
}

func CreateNGramSet(n int) *NGramSet {
	return &NGramSet{n, 0, make(map[string]int64)}
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
		universe.NgramsByN[n].Print()
		n++
	}
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
	set := NGramSet{n, 0, make(map[string]int64)}

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
