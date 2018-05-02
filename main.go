package main

import (
	"fmt"
)

func doNothing(state ExperimentState) []Order {
	var orders []Order
	return orders
}

func main() {
	companiesBySymbol := GetCompaniesBySmybol()

	var symbols []string
	for _, company := range companiesBySymbol {
		symbols = append(symbols, company.Symbol)
	}

	fmt.Println(symbols)
	GetDailyStockSummaryData()
	fmt.Println("Loaded summary data")

	/*
	for _, summary := range GetDailyStockSummaryData().OnDay(ParseDay("2017-11-01")) {
		fmt.Println(summary)
	}

	return
	*/
	// allData := GetAllDailyStockPrices()

	// fmt.Println(len(allData))

	params := ExperimentParams{}
	params.InitialBalance = 1000
	params.TransactionFee = 12
	params.StartDay = ParseDay("2017-01-01")
	params.EndDay = ParseDay("2018-01-01")
	params.CompaniesBySymbol = companiesBySymbol

	experiment := CreateExperiment(params, OptimalDayTradingStategy)
	experiment.Run()
}
