package main

import (
	"fmt"
)

func getGradeString(quotes []*Quote) string {
	gradeString := ""

	for _, quote := range quotes {
		grade := quote.GetGrade()
		gradeString = gradeString + grade
	}

	return gradeString
}



func PrintPredictionsForNextWeek() {
	maxN := 3
	db := GetHandle()

	tickers := GetAllCompanyTickers(db)[:1000]

	ngramModel := CreateNGramModel(GradeQuoteClassifier, "", maxN)

	for i, ticker := range tickers {
		fmt.Printf("[TRAIN] %4d of %4d   (%s)\n", i, len(tickers), ticker)
		quotes := GetAllQuotesForTicker(db, ticker)

		if len(quotes) < maxN {
			continue
		}

		company := GetCompanyByTicker(db, ticker)
		document := getGradeString(quotes)
		ngramModel.AddCase(company, document)
	}

	fmt.Printf("\n\n")

	for _, ticker := range tickers {
		company := GetCompanyByTicker(db, ticker)
		quotes := GetLatestQuotesForTicker(db, ticker, maxN)

		if len(quotes) < maxN {
			continue
		}

		recentGrades := getGradeString(quotes)
		predictedGrade, probability := ngramModel.PredictNext(company, recentGrades)
		latestQuote := quotes[len(quotes) - 1]

		fmt.Printf("[PREDICT] %6s @ %s \t\t|\t\t%s -> %s (%2.2f%%)\n",
			ticker,
			TimeToString(latestQuote.Date),
			recentGrades,
			predictedGrade,
			probability*100)
		// ngramModel.AddCase(company, document)
	}
}