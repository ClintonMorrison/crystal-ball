package main

import (
	"fmt"
	"strings"
)


func PrintPredictionsForNextWeek() {
	maxN := 3
	db := GetHandle()

	tickers := GetAllCompanyTickers(db)

	ngramModel := CreateNGramModel(GradeQuoteClassifier, "", maxN)

	for i, ticker := range tickers {
		fmt.Printf("[TRAIN] %4d of %4d   (%s)\n", i, len(tickers), ticker)
		quotes := GetAllQuotesForTicker(db, ticker)

		if len(quotes) < maxN {
			continue
		}

		company := GetCompanyByTicker(db, ticker)
		document := GetGradeString(quotes)
		ngramModel.AddCase(company, document)
	}

	fmt.Printf("\n\n")

	bestTickers := make([]string, 0)
	bestProbability := 0.0


	for _, ticker := range tickers {
		company := GetCompanyByTicker(db, ticker)
		quotes := GetLatestQuotesForTicker(db, ticker, maxN)

		if len(quotes) < maxN {
			continue
		}

		recentGrades := GetGradeString(quotes)
		predictedGrade, probability := ngramModel.PredictNext(company, recentGrades)
		latestQuote := quotes[len(quotes) - 1]

		if predictedGrade == "A" {
			if probability > bestProbability {
				bestTickers = make([]string, 0)
				bestTickers = append(bestTickers, ticker)
				bestProbability = probability
			} else if probability == bestProbability {
				bestTickers = append(bestTickers, ticker)
			}

		}

		fmt.Printf("[PREDICT] %6s @ %s \t\t|\t\t%s -> %s (%2.2f%%)\n",
			ticker,
			TimeToString(latestQuote.Date),
			recentGrades,
			predictedGrade,
			probability*100)
	}

	fmt.Printf("\n\nBest probability: %2.2f%%\n", bestProbability * 100)
	fmt.Printf("Tickers: %s\n", strings.Join(bestTickers, ", "))
}