package main

import (
	"stock-analysis/experiment"
	"stock-analysis/util"
	"stock-analysis/data"
	"stock-analysis/stategies"
)

func runExperiment(companiesBySymbol map[string]data.Company) {
	params := experiment.ExperimentParams{}
	params.InitialBalance = 1000.0
	params.TransactionFee = 0 // 12
	params.StartDay = util.ParseDay("2017-01-01")
	params.EndDay = util.ParseDay("2018-01-01")
	params.CompaniesBySymbol = companiesBySymbol

	experiment := experiment.CreateExperiment(params, stategies.MovingAverageTrendFollowing)
	experiment.Run()
}

func main() {
	data.GetDailyStockSummaryData()
	companiesBySymbol := data.GetCompaniesBySmybol()
	runExperiment(companiesBySymbol)
}
