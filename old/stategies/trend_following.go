package stategies

import (
	"stockAnalysis/experiment"
	"stockAnalysis/data"
)

func MovingAverageTrendFollowing(state *experiment.ExperimentState) []experiment.Order {
	var orders []experiment.Order

	// Sell all owned stocks, that are trending down
	for symbol, qty := range state.Portfolio {
		stock := data.GetDailyStockSummaryData().ForSymbolOnDay(symbol, state.Day)
		prediction := PredictDirectionWithMovingAverage(stock)
		if prediction == "DOWN" && qty > 0 {
			orders = append(orders, experiment.Order{
				experiment.SELL,
				symbol,
				qty,
			})
		}
	}

	stocks := data.GetDailyStockSummaryData().OnDay(state.Day)
	for _, stock := range stocks {
		prediction := PredictDirectionWithMovingAverage(stock)

		if prediction == "UP" && state.Balance > 10.0 {
			orders = append(orders, experiment.Order{
				experiment.BUY,
				stock.Symbol,
				state.Balance / stock.Close })
		}
	}

	return orders
}

func PredictDirectionWithMovingAverage(summary data.StockSummary) string {
	if summary.MovingAverage(2) > summary.MovingAverage(4) {
		return "UP"
	} else if summary.MovingAverage(2) < summary.MovingAverage(4) {
		return "DOWN"
	}

	return "STABLE"
}
