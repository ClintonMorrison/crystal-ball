package main

import (
	"fmt"
	"time"
	"stock-analysis/trading"
	"math"
)

func parseDay(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func doNothing(state trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	return orders
}

func buyAll(state *trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	orders = append(orders, trading.Order{trading.BUY, "CZFC", 1.0})
	allStocks := trading.GetAvailableStocksForDay(state.Day)
	fmt.Println(len(allStocks))
	return orders
}

func scoreStock(stock trading.Stock) float64 {
	percentChange := (stock.Close - stock.Open) / stock.Close
	weightedChange := percentChange * 100.0 * math.Pow(stock.Volume, 0.1)
	return weightedChange
}


func tradeByScore(state *trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	stocksBySymbol := trading.GetAvailableStocksBySymbol(state.Day)

	maxScoreStock := trading.Stock{}
	maxScore := 0.0

	minScoreStock := trading.Stock{}
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
		orders = append(orders, trading.Order{trading.BUY, maxScoreStock.Symbol, 100.0 / maxScoreStock.Close })
	}

	if len(minScoreStock.Symbol) > 0 {
		fmt.Printf("min score was %.2f [$%.2f --> $%.2f]\n", minScore, minScoreStock.Open, minScoreStock.Close)
		orders = append(orders, trading.Order{trading.SELL, minScoreStock.Symbol, state.Portfolio[minScoreStock.Symbol] })
	}

	return orders
}


func buyBiggestDailyChange(state *trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	allStocks := trading.GetAvailableStocksForDay(state.Day)

	maxIncreaseStock := trading.Stock{}
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
		orders = append(orders, trading.Order{trading.BUY, maxIncreaseStock.Symbol, 10.0})
	}

	// orders = append(orders, trading.Order{trading.BUY, "CZFC", 1.0})
	return orders
}

func main() {
	companiesBySymbol := trading.GetCompaniesBySmybol()

	var symbols []string
	for _, company := range companiesBySymbol {
		symbols = append(symbols, company.Symbol)
	}

	fmt.Println(symbols)

	allData := trading.GetStocksDailyData()
	
	fmt.Println(len(allData))
	/*

	params := trading.ExperimentParams{}
	params.InitialBalance = 1000
	params.StartDay = parseDay("2015-11-01")
	params.EndDay = parseDay("2018-02-01")
	params.CompaniesBySymbol = companiesBySymbol
	trading.RunExperiment(params, tradeByScore)
	*/
}
