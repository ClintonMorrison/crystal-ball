package main

import (
	"fmt"
)

func EvaluateNGramModelOverTime () {
	maxN := 3
	splitDate := ParseDateString("2016-01-01")
	endDate := ParseDateString("2018-10-28")
	company := Company{}

	db := GetHandle()

	tickers := GetAllCompanyTickers(db)[:1]
	dates := make([]string, 0)

	fmt.Print(tickers)

	fmt.Printf("\n\n\n%10s\t%10s\t%s\n", "Index", "Date", "Performance")

	i := 1

	for splitDate.Before(endDate) {
		dates = append(dates, TimeToString(splitDate))
		splitDate = AddDays(splitDate, 28)

		trainingDocuments := make(map[string]string, 0)
		evaluationDocuments := make(map[string]string, 0)

		for _, ticker := range tickers {
			quotes := GetAllQuotesForTicker(db, ticker)
			gradeString, trainingString := getSplitGradeString(quotes, splitDate)
			trainingDocuments[ticker] = trainingString
			evaluationDocuments[ticker] = gradeString
		}

		ngramModel := CreateNGramModel(GradeQuoteClassifier, "", maxN)
		evaluationResults := EvaluationResults{}
		evaluationResults.init()

		for _, document := range trainingDocuments {
			ngramModel.AddCase(&company, document)
		}

		for ticker, document := range evaluationDocuments {
			trainingDocument := trainingDocuments[ticker]
			trainingEnd := len(trainingDocument)
			nextChar := trainingEnd + 1
			if len(trainingDocument) < maxN || len(document) < len(trainingDocument) + 1  {
				// fmt.Printf("[WARN] skipping %s because not enough training data\n", ticker)
				continue
			}

			pred, _ := ngramModel.PredictNext(&company, trainingDocument)
			actual := document[trainingEnd:nextChar]
			evaluationResults.AddCase(actual, pred)
		}

		fmt.Printf("%10d\t%10s\t%2.2f%%\n", i, TimeToString(splitDate), evaluationResults.PercentCorrect())
		i++
	}
}
