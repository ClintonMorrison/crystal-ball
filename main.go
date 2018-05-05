package main

import (
	"fmt"
	"stock-analysis/experiment"
	"stock-analysis/util"
	"stock-analysis/data"
	"stock-analysis/stategies"
)

func doNothing(state experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order
	return orders
}

func runExperiment(companiesBySymbol map[string]data.Company) {
	params := experiment.ExperimentParams{}
	params.InitialBalance = 1000
	params.TransactionFee = 0 // 12
	params.StartDay = util.ParseDay("2017-01-01")
	params.EndDay = util.ParseDay("2018-01-01")
	params.CompaniesBySymbol = companiesBySymbol

	experiment := experiment.CreateExperiment(params, stategies.KnnByDistance) // OptimalDayTradingStategy)
	experiment.Run()
}

func main() {
	companiesBySymbol := data.GetCompaniesBySmybol()

	var symbols []string
	for _, company := range companiesBySymbol {
		symbols = append(symbols, company.Symbol)
	}

	data.GetDailyStockSummaryData()

	totalExamples := 0
	numberCorrect := 0
	for _, stocksOnDay := range data.GetDailyStockSummaryData().SummariesByDay {
		for _, stock := range stocksOnDay {
			predictedChange := stategies.PredictDirectionWithNearest(stock)
			nextDaySummary := data.GetDailyStockSummaryData().ForSymbolOnOtherDay(stock.Symbol, stock.Date, 1)
			actualChange := "DOWN"

			if nextDaySummary.GetPercentChange() > 0 {
				actualChange = "UP"
			}

			totalExamples++
			if actualChange == predictedChange {
				numberCorrect++
			}

			// fmt.Printf("%s [%s] --> %s / %s [%s]\n", TimeToString(day), stock.Symbol, predictedChange, actualChange, nextDaySummary.Symbol)
		}

		if totalExamples > 5000 {
			break
		}
	}

	fmt.Printf("\nTotal: %d\nCorrect: %d\nAccuracy: %%%.2f", totalExamples, numberCorrect, float64(numberCorrect)/float64(totalExamples)*100.0)

	// runExperiment(companiesBySymbol)
}
