package main

import (
	"fmt"
)

func EvaluateNGramModelOverTime () {
	maxN := 3
	splitDate := ParseDateString("2000-01-01")
	endDate := ParseDateString("2018-10-28")

	resultsByConfiguration := make(map[string]*EvaluationResults, 0)

	db := GetHandle()

	tickers := GetAllCompanyTickers(db)
	dates := make([]string, 0)

	fmt.Print(tickers)

	companies := make(map[string]*Company, 0)
	for _, ticker := range tickers {
		company := GetCompanyByTicker(db, ticker)
		companies[ticker] = company
	}


	for splitDate.Before(endDate) {
		dates = append(dates, TimeToString(splitDate))
		splitDate = AddDays(splitDate, 7)
		fmt.Println(TimeToString(splitDate))
		fmt.Printf("Building model for %s\n", TimeToString(splitDate))

		trainingDocuments := make(map[string]string, 0)
		evaluationDocuments := make(map[string]string, 0)

		for _, ticker := range tickers {
			quotes := GetAllQuotesForTicker(db, ticker)
			gradeString, trainingString := getSplitGradeString(quotes, splitDate)
			trainingDocuments[ticker] = trainingString
			evaluationDocuments[ticker] = gradeString
		}

		ngramModel := CreateNGramModel(GradeQuoteClassifier, "", maxN)
		key := TimeToString(splitDate)
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

	i := 0
	fmt.Printf("\n\n\n%10s\t%10s\t%s\n", "Index", "Date", "Performance")
	for _, date := range dates {
		results := resultsByConfiguration[date]
		if results != nil {
			i++
			fmt.Printf("%10d\t%10s\t%2.2f%%\n", i, date, results.PercentCorrect())
		}
	}
}
