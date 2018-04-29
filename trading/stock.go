package trading

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Stock struct {
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

func StocksByDay(stocks []Stock) map[string]Stock {
	stocksByDay := make(map[string]Stock)
	for _, stock := range stocks {
		stocksByDay[TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

type DailyStocksBySymbol map[string]map[string]Stock
