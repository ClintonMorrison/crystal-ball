package main

import (
	"fmt"
	"math"
)



func doNothing(state ExperimentState) []Order {
	var orders []Order
	return orders
}

func buyAll(state *ExperimentState) []Order {
	var orders []Order
	orders = append(orders, Order{BUY, "CZFC", 1.0})
	allStocks := GetStockSummariesForDay(state.Day)
	fmt.Println(len(allStocks))
	return orders
}

func scoreStock(stock StockSummary) float64 {
	percentChange := (stock.Close - stock.Open) / stock.Close
	weightedChange := percentChange * 100.0 * math.Pow(stock.Volume, 0.1)
	return weightedChange
}


func tradeByScore(state *ExperimentState) []Order {
	var orders []Order
	stocksBySymbol := GetAvailableStocksBySymbol(state.Day)

	maxScoreStock := StockSummary{}
	maxScore := 0.0

	minScoreStock := StockSummary{}
	minScore := 0.0

	for _, stock := range stocksBySymbol {
		score := scoreStock(stock)
		owned := state.Portfolio[stock.Symbol] > 0

		if score > maxScore {
			maxScoreStock = stock
			maxScore = score
		}

		lessThanMin := score < minScore
		if owned && lessThanMin {
			minScoreStock = stock
			minScore = score
		}
	}

	if len(maxScoreStock.Symbol) > 0 {
		fmt.Printf("max score was %.2f [$%.2f --> $%.2f]\n", maxScore, maxScoreStock.Open, maxScoreStock.Close)
		orders = append(orders, Order{BUY, maxScoreStock.Symbol, 100.0 / maxScoreStock.Close })
	}

	if len(minScoreStock.Symbol) > 0 {
		fmt.Printf("min score was %.2f [$%.2f --> $%.2f]\n", minScore, minScoreStock.Open, minScoreStock.Close)
		orders = append(orders, Order{SELL, minScoreStock.Symbol, state.Portfolio[minScoreStock.Symbol] })
	}

	return orders
}


func buyBiggestDailyChange(state *ExperimentState) []Order {
	var orders []Order
	allStocks := GetStockSummariesForDay(state.Day)

	maxIncreaseStock := StockSummary{}
	maxIncrease := 0.0

	for _, stock := range allStocks {
		change := (stock.Close - stock.Open) / stock.Close
		if change < maxIncrease {
			maxIncreaseStock = stock
			maxIncrease = change
		}
	}

	if len(maxIncreaseStock.Symbol) > 0 {
		fmt.Printf("max increase was %%%.2f [$%.2f --> $%.2f]\n", maxIncrease * 100, maxIncreaseStock.Open, maxIncreaseStock.Close)
		orders = append(orders, Order{BUY, maxIncreaseStock.Symbol, 10.0})
	}

	// orders = append(orders, Order{BUY, "CZFC", 1.0})
	return orders
}

func main() {
	companiesBySymbol := GetCompaniesBySmybol()

	var symbols []string
	for _, company := range companiesBySymbol {
		symbols = append(symbols, company.Symbol)
	}

	fmt.Println(symbols)

	// allData := GetAllDailyStockPrices()
	
	// fmt.Println(len(allData))

	params := ExperimentParams{}
	params.InitialBalance = 1000
	params.StartDay = ParseDay("2015-11-01")
	params.EndDay = ParseDay("2018-02-01")
	params.CompaniesBySymbol = companiesBySymbol
	experiment := CreateExperiment(params, tradeByScore)
	experiment.Run()
}
