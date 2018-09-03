package stategies

import (
	"sort"
	"time"
	"stock-analysis/data"
	"stock-analysis/util"
	"stock-analysis/experiment"
)

func NearestPreviousSummaries(summary data.StockSummary, k int) []data.StockSummary {
	previousSummaries := data.GetDailyStockSummaryData().InDateRange(
		summary.Date.Add(time.Hour * time.Duration(-10 * 24)),
		summary.Date.Add(time.Hour * time.Duration(-1 * 24)))

	return findNearestNeighbors(k, summary, previousSummaries)
}

func PredictChangeWithNearest(summary data.StockSummary) float64 {
	neighbors := NearestPreviousSummaries(summary, 10)
	var changes []float64
	for _, neighbor := range neighbors {
		nextDayForNeighbor := data.GetDailyStockSummaryData().ForSymbolOnOtherDay(neighbor.Symbol, neighbor.Date, 1)
		if nextDayForNeighbor.Symbol != "" {
			changes = append(changes, nextDayForNeighbor.GetPercentChange())
		}
	}

	if len(changes) == 0 {
		return 0.0
	}

	return util.Avg(changes)
}


const (
	UP = iota
	DOWN
)

func PredictDirectionWithNearest(summary data.StockSummary) string {
	neighbors := NearestPreviousSummaries(summary, 10)
	var changes []float64
	for _, neighbor := range neighbors {
		nextDayForNeighbor := data.GetDailyStockSummaryData().ForSymbolOnOtherDay(neighbor.Symbol, neighbor.Date, 1)
		if nextDayForNeighbor.Symbol != "" {
			changes = append(changes, nextDayForNeighbor.GetPercentChange())
		}
	}

	numberPositive := util.CountAboveThreshold(0, changes)
	midpoint := len(neighbors) / 2

	if numberPositive > midpoint {
		return "UP"
	}

	return "DOWN"
}

type byDistance []data.StockSummary
var subjectStock = data.StockSummary{}

func (s byDistance) Len() int {
	return len(s)
}
func (s byDistance) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDistance) Less(i, j int) bool {
	return s[i].Distance(subjectStock) < s[j].Distance(subjectStock)
}

func findNearestNeighbors(k int, summary data.StockSummary, summaries []data.StockSummary) []data.StockSummary{
	subjectStock = summary
	sort.Sort(byDistance(summaries))
	return summaries[:k]
}

func KnnByDistance(state *experiment.ExperimentState) []experiment.Order {
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

	yesterday := state.Day.Add(-24 * time.Hour)
	stocksForYesterday := data.GetDailyStockSummaryData().OnDay(yesterday)
	stocksForToday := data.GetDailyStockSummaryData().OnDay(state.Day)

	noStocksYesterday := len(stocksForToday) == 0
	noStocksToday := len(stocksForYesterday) == 0
	if noStocksYesterday || noStocksToday {
		return orders
	}

	stockWithGreatestExpectedReturn := data.StockSummary{}
	expectedPercentChange := 0.0

	for _, stock := range stocksForToday {
		neighbors := NearestPreviousSummaries(stock, 5)

		var neighborPercentChanges []float64

		for _, neighbor := range neighbors {
			neighborToday := data.GetDailyStockSummaryData().ForSymbolOnDay(neighbor.Symbol, state.Day)
			neighborPercentChanges = append(neighborPercentChanges, neighborToday.GetPercentChange())
		}
		neighborAverage := util.Avg(neighborPercentChanges)
		numberAboveThreshold := util.CountAboveThreshold(0.01, neighborPercentChanges)



		if neighborAverage > expectedPercentChange && numberAboveThreshold >= 5 {
			stockWithGreatestExpectedReturn = stock
			expectedPercentChange = neighborAverage
		}
	}

	// Buy best looking stock
	if stockWithGreatestExpectedReturn.Symbol != "" {
		orders = append(orders, experiment.Order{
			experiment.BUY,
			stockWithGreatestExpectedReturn.Symbol,
			state.Balance / stockWithGreatestExpectedReturn.Close,
		})
	}

	return orders
}
