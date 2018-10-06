package main

import (
	"flag"
)

func main() {

	evaluateModels := flag.Bool(
		"evaluate-ngram-models",
		false,
		"Evaluate different parameters for building n-gram models")

	predictNextWeek := flag.Bool(
		"predict",
		false,
		"Predict stock changes for the next week")

	downloadWeeklyQuotes := flag.Bool(
		"download-quotes",
		false,
		"Download weekly stock price data")

	downloadCompanies := flag.Bool(
		"download-companies",
		false,
		"Download list of current companies")

	flag.Parse()

	if *predictNextWeek {
		PrintPredictionsForNextWeek()
	} else if *evaluateModels {
		EvaluateNGramModels()
	} else if *downloadWeeklyQuotes {
		DownloadWeeklyQuotes()
	} else if *downloadCompanies {
		DownloadCompanies()
	}

	flag.Usage()


	// Print Stock Price

	// Download Companies

	// Download Quotes

}
