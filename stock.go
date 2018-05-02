package main

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type StockSummary struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Symbol           string
	Date             time.Time
	Open             float64
	High             float64
	Low              float64
	Close            float64
	AdjustedClose    float64
	Volume           float64
	DividendAmount   float64
	SplitCoefficient float64
}


func (summary StockSummary) getChange() float64 {
	return summary.Close - summary.Open
}

func (summary StockSummary) getPercentChange() float64 {
	return (summary.Close - summary.Open) / summary.Open
}


func GroupStocksByDay(stocks []StockSummary) map[string]StockSummary {
	stocksByDay := make(map[string]StockSummary)
	for _, stock := range stocks {
		stocksByDay[TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

func FindStockWithMaximumIncrease(summaries []StockSummary) StockSummary {
	maxChange := 0.0
	maxSummary := StockSummary{}

	for _, summary := range summaries {
		change := summary.getPercentChange()
		if change > maxChange {
			maxChange = change
			maxSummary = summary
		}
	}

	return maxSummary
}

func GetAvailableStocksBySymbol(date time.Time) map[string]StockSummary {
	stocks := GetStockSummariesForDay(date)
	stocksBySmybol := make(map[string]StockSummary)
	for _, stock := range stocks {
		stocksBySmybol[stock.Symbol] = stock
	}

	return stocksBySmybol
}

func GetPricesForStockByDay(symbol string) map[string]StockSummary {
	stocks := GetDailySummariesForStock(symbol)
	stocksByDay := make(map[string]StockSummary, len(stocks))
	for _, stock := range stocks {
		stocksByDay[TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

func (collection SummaryData) forSymbolOnDay(symbol string, date time.Time) StockSummary {
	return collection.summariesBySymbolByDay[symbol][date]
}

func (collection SummaryData) OnDay(date time.Time) []StockSummary {
	return collection.summariesByDay[date]
}

var cachedSummaryData SummaryData
var generatedSummaryData = false

func GetDailyStockSummaryData() SummaryData {
	if generatedSummaryData {
		return cachedSummaryData
	}

	summaries := GetAllDailyStockPrices()

	summariesByDay := make(map[time.Time][]StockSummary)
	summariesBySymbol := make(map[string][]StockSummary)
	summariesBySymbolByDay := make(map[string]map[time.Time]StockSummary)

	for _, summary := range summaries {
		summariesByDay[summary.Date] = append(summariesByDay[summary.Date], summary)
		summariesBySymbol[summary.Symbol] = append(summariesBySymbol[summary.Symbol], summary)

		if summariesBySymbolByDay[summary.Symbol] == nil {
			summariesBySymbolByDay[summary.Symbol] = make(map[time.Time]StockSummary)
		}
		summariesBySymbolByDay[summary.Symbol][summary.Date] = summary
	}

	summaryData := SummaryData{
		summariesByDay,
		summariesBySymbol,
		summariesBySymbolByDay,
	}
	cachedSummaryData = summaryData
	generatedSummaryData = true
	return summaryData
}

type SummaryData struct {
	summariesByDay         map[time.Time][]StockSummary
	summariesBySymbol      map[string][]StockSummary
	summariesBySymbolByDay map[string]map[time.Time]StockSummary
}

type DailyStocksBySymbol map[string]map[string]StockSummary
