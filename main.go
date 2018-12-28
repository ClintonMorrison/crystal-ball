package main

import (
	"flag"
)

func main() {

	evaluateModelsParams := flag.Bool(
		"evaluate-ngram-model-params",
		false,
		"Evaluate different parameters for building n-gram models")

	evaluateNgramModelOverTime := flag.Bool(
		"evaluate-ngram-model-over-time",
		false,
		"Evaluates an ngram model with n=3 over time")

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
		"Flag specifying ticker to start at (for -download-quotes)")

	flag.Parse()

	if *predictNextWeek {
		PrintPredictionsForNextWeek()
	} else if *evaluateModelsParams {
		EvaluateNGramModels()
	} else if *downloadWeeklyQuotes {
		DownloadWeeklyQuotes(*startTicker)
	} else if *downloadCompanies {
		DownloadCompanies()
	} else if *printNgrams {
		PrintNgrams()
	} else if *evaluateNgramModelOverTime {
		EvaluateNGramModelOverTime()
	} else {
		flag.Usage()
	}
}
