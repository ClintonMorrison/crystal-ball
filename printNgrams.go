package main

import "fmt"

func PrintNgrams() {
	maxN := 3
	db := GetHandle()

	tickers := GetAllCompanyTickers(db)

	universe := CreateUniverse(maxN)

	for i, ticker := range tickers {
		fmt.Printf("[TRAIN] %4d of %4d   (%s)\n", i, len(tickers), ticker)
		quotes := GetAllQuotesForTicker(db, ticker)

		if len(quotes) < maxN {
			continue
		}

		document := GetGradeString(quotes)
		universe.AddString(document)
	}

	universe.Print()
}