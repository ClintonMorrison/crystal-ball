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

	printNgrams := flag.Bool(
		"print-ngrams",
		false,
		"Prints n-gram frequencies for aggregate stock data")

	startTicker := flag.String(
		"start-ticker",
		"",
		"Ticker to start at")

	flag.Parse()

	if *predictNextWeek {
		PrintPredictionsForNextWeek()
	} else if *evaluateModels {
		EvaluateNGramModels()
	} else if *downloadWeeklyQuotes {
		DownloadWeeklyQuotes(*startTicker)
	} else if *downloadCompanies {
		DownloadCompanies()
	} else if *printNgrams {
		PrintNgrams()
	} else {
		flag.Usage()
	}
}
