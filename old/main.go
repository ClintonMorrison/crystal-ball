package main

import (
	"stock-analysis/experiment"
	"stock-analysis/util"
	"stock-analysis/data"
	"stock-analysis/stategies"
	"fmt"
)

func runExperiment(companiesBySymbol map[string]data.Company) {
	params := experiment.ExperimentParams{}
	params.InitialBalance = 1000.0
	params.TransactionFee = 0 // 12
	params.StartDay = util.ParseDay("2012-01-01")
	params.EndDay = util.ParseDay("2018-01-01")
	params.CompaniesBySymbol = companiesBySymbol

	experiment := experiment.CreateExperiment(params, stategies.MovingAverageTrendFollowing) // stategies.MovingAverageTrendFollowing)
	experiment.Run()
}

func main() {
	data.GetDailyStockSummaryData()
	fmt.Print(data.GetDailyStockSummaryData().ForSymbolOnDay("GOOG", util.ParseDay("2017-08-17")))
	// companiesBySymbol := data.GetCompaniesBySmybol()
	// runExperiment(companiesBySymbol)
}
