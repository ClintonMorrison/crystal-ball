package stategies

import (
	"time"
	"stock-analysis/experiment"
	"stock-analysis/data"
)

func OptimalDayTradingStategy(state *experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order

	// Sell all owned stocks
	for symbol, qty := range state.Portfolio {
		if qty > 0 {
			orders = append(orders, experiment.Order{
				experiment.SELL,
				symbol,
				qty,
			})
		}
	}

	// Buy tomorrow's biggest increase
	tomorrow := state.Day.Add(24 * time.Hour)
	summaries := data.GetDailyStockSummaryData().OnDay(tomorrow)

	if len(summaries) > 0 {
		maxIncreaseStock := data.FindStockWithMaximumIncrease(summaries)

		qtyToBuy := state.Balance / maxIncreaseStock.Open

		orders = append(orders, experiment.Order{
			experiment.BUY,
			maxIncreaseStock.Symbol,
			qtyToBuy,
		})
	}

	return orders
}