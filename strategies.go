package main

import (
	"fmt"
	"math"
	"time"
)

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
		orders = append(orders, Order{BUY, maxScoreStock.Symbol, 100.0 / maxScoreStock.Close})
	}

	if len(minScoreStock.Symbol) > 0 {
		fmt.Printf("min score was %.2f [$%.2f --> $%.2f]\n", minScore, minScoreStock.Open, minScoreStock.Close)
		orders = append(orders, Order{SELL, minScoreStock.Symbol, state.Portfolio[minScoreStock.Symbol]})
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
		fmt.Printf("max increase was %%%.2f [$%.2f --> $%.2f]\n", maxIncrease*100, maxIncreaseStock.Open, maxIncreaseStock.Close)
		orders = append(orders, Order{BUY, maxIncreaseStock.Symbol, 10.0})
	}

	// orders = append(orders, Order{BUY, "CZFC", 1.0})
	return orders
}

func OptimalDayTradingStategy(state *ExperimentState) []Order {
	var orders []Order

	// Sell all owned stocks
	for symbol, qty := range state.Portfolio {
		if qty > 0 {
			orders = append(orders, Order{
				SELL,
				symbol,
				qty,
			})
		}
	}

	// Buy tomorrow's biggest increase
	tomorrow := state.Day.Add(24 * time.Hour)
	summaries := GetDailyStockSummaryData().OnDay(tomorrow)

	if len(summaries) > 0 {
		maxIncreaseStock := FindStockWithMaximumIncrease(summaries)

		qtyToBuy := state.Balance / maxIncreaseStock.Open

		orders = append(orders, Order{
			BUY,
			maxIncreaseStock.Symbol,
			qtyToBuy,
		})
	}

	return orders
}
