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

func StocksByDay(stocks []StockSummary) map[string]StockSummary {
	stocksByDay := make(map[string]StockSummary)
	for _, stock := range stocks {
		stocksByDay[TimeToString(stock.Date)] = stock
	}

	return stocksByDay
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


type DailyStocksBySymbol map[string]map[string]StockSummary
