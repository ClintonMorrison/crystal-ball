package main

import (
	"fmt"
	"time"
)


func getSplitGradeString(quotes []*Quote, splitTime time.Time) (string, string) {
	gradeString := ""
	trainingString := ""

	for _, quote := range quotes {
		grade := quote.GetGrade()
		gradeString = gradeString + grade

		if quote.Date.Before(splitTime) {
			trainingString = trainingString + grade
		}
	}

	return gradeString, trainingString
}

func EvaluateNGramModels () {
	params := []string{ "INDUSTRY", "SECTOR", "TICKER", "" }
	nValues := []int { 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18 }
	splitDate := ParseDateString("2018-01-01")

	resultsByConfiguration := make(map[string]*EvaluationResults, 0)

	db := GetHandle()

	tickers := GetAllCompanyTickers(db)

	trainingDocuments := make(map[string]string, 0)
	evaluationDocuments := make(map[string]string, 0)
	companies := make(map[string]*Company, 0)

	for _, ticker := range tickers {
		quotes := GetAllQuotesForTicker(db, ticker)
		company := GetCompanyByTicker(db, ticker)

		companies[ticker] = company
		gradeString, trainingString := getSplitGradeString(quotes, splitDate)
		trainingDocuments[ticker] = trainingString
		evaluationDocuments[ticker] = gradeString
	}

	for _, parameter := range params {
		for _, maxN := range nValues {
			fmt.Printf("%s %d \n", parameter, maxN)
			ngramModel := CreateNGramModel(GradeQuoteClassifier, parameter, maxN)
			key := parameter + string(maxN)
			evaluationResults := EvaluationResults{}
			evaluationResults.init()
			resultsByConfiguration[key] = &evaluationResults

			for ticker, document := range trainingDocuments {
				ngramModel.AddCase(companies[ticker], document)
			}

			for ticker, document := range evaluationDocuments {
				company := companies[ticker]
				trainingDocument := trainingDocuments[ticker]
				trainingEnd := len(trainingDocument)
				nextChar := trainingEnd + 1
				if len(trainingDocument) < maxN || len(document) < len(trainingDocument) + 1  {
					fmt.Printf("[WARN] skipping %s because not enough training data\n", ticker)
					continue
				}

				pred, _ := ngramModel.PredictNext(company, trainingDocument)
				actual := document[trainingEnd:nextChar]
				evaluationResults.AddCase(actual, pred)
			}
		}
	}

	fmt.Printf("\n\n\n%10s\t%4s\t%s\n", "Parameter", "MaxN", "Performance")
	for _, parameter := range params {
		for _, maxN := range nValues {
			key := parameter + string(maxN)
			results := resultsByConfiguration[key]
			fmt.Printf("%10s\t%4d\t%2.2f%%\n", parameter, maxN, results.PercentCorrect())
		}
	}
}
