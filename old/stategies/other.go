package stategies

import (
	"math"
	"fmt"
	"stockAnalysis/data"
	"stockAnalysis/experiment"
)

func scoreStock(stock data.StockSummary) float64 {
	percentChange := (stock.Close - stock.Open) / stock.Close
	weightedChange := percentChange * 100.0 * math.Pow(stock.Volume, 0.1)
	return weightedChange
}

func tradeByScore(state *experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order
	stocksBySymbol := data.GetAvailableStocksBySymbol(state.Day)

	maxScoreStock := data.StockSummary{}
	maxScore := 0.0

	minScoreStock := data.StockSummary{}
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
		orders = append(orders, experiment.Order{experiment.BUY, maxScoreStock.Symbol, 100.0 / maxScoreStock.Close})
	}

	if len(minScoreStock.Symbol) > 0 {
		fmt.Printf("min score was %.2f [$%.2f --> $%.2f]\n", minScore, minScoreStock.Open, minScoreStock.Close)
		orders = append(orders, experiment.Order{experiment.SELL, minScoreStock.Symbol, state.Portfolio[minScoreStock.Symbol]})
	}

	return orders
}

func buyBiggestDailyChange(state *experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order
	allStocks := data.GetStockSummariesForDay(state.Day)

	maxIncreaseStock := data.StockSummary{}
	maxIncrease := 0.0

	for _, stock := range allStocks {
		change := (stock.Close - stock.Open) / stock.Close
		if change < maxIncrease {
			maxIncreaseStock = stock
			maxIncrease = change
		}
	}

	if len(maxIncreaseStock.Symbol) > 0 {
		fmt.Printf("max increase was %%%.2f [$%.2f --> $%.2f]\n", maxIncrease*100, maxIncreaseStock.Open, maxIncreaseStock.Close)
		orders = append(orders, experiment.Order{experiment.BUY, maxIncreaseStock.Symbol, 10.0})
	}

	// orders = append(orders, Order{BUY, "CZFC", 1.0})
	return orders
}


func buyAll(state *experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order
	orders = append(orders, experiment.Order{experiment.BUY, "CZFC", 1.0})
	allStocks := data.GetStockSummariesForDay(state.Day)
	fmt.Println(len(allStocks))
	return orders
}
