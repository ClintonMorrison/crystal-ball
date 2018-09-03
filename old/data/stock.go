package data

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"fmt"
	"math"
	"stock-analysis/util"
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

func (summary StockSummary) Print() {
	fmt.Println(summary.Symbol)
	fmt.Println("-------")
	fmt.Printf("Date: %s\n", util.TimeToString(summary.Date))
	fmt.Printf("Open: $%.2f\n", summary.Open)
	fmt.Printf("Close: $%.2f\n", summary.Close)
	fmt.Printf("Low: $%.2f\n", summary.Low)
	fmt.Printf("High: $%.2f\n", summary.High)
	fmt.Printf("Volume: %.2f\n", summary.Volume)
	fmt.Printf("Change: %%%.2f\n", summary.GetPercentChange() * 100.0)
	fmt.Printf("Violatility: %%%.2f\n", summary.GetPercentVolatility() * 100.0)

	fmt.Println()
}

func (summary StockSummary) GetChange() float64 {
	return summary.Close - summary.Open
}

func (summary StockSummary) GetPercentChange() float64 {
	return (summary.Close - summary.Open) / summary.Open
}

func (summary StockSummary) GetPercentVolatility() float64 {
	return (summary.High - summary.Low) / summary.Open
}

func (summary StockSummary) GetVelocity(days int) float64 {
	summaries := GetDailyStockSummaryData().ForSymbolOnOtherDays(summary.Symbol, summary.Date, -days)

	totalChange := 0.0

	for _, summary := range summaries {
		totalChange += summary.GetChange()
	}

	return totalChange / float64(len(summaries))
}

func (summary StockSummary) MovingAverage(days int) float64 {
	summaries := GetDailyStockSummaryData().ForSymbolOnOtherDays(summary.Symbol, summary.Date, -days)

	var closeValues []float64
	for _, summary := range summaries {
		closeValues = append(closeValues, summary.GetChange())
	}

	return util.Avg(closeValues)
}

func (s1 StockSummary) Distance(s2 StockSummary) float64 {
	x1 := s2.GetPercentChange() - s1.GetPercentChange()
	x2 := s2.GetPercentVolatility() - s1.GetPercentVolatility()
	return math.Sqrt(math.Pow(x1, 2) + math.Pow(x2, 2))
}



func GroupStocksByDay(stocks []StockSummary) map[string]StockSummary {
	stocksByDay := make(map[string]StockSummary)
	for _, stock := range stocks {
		stocksByDay[util.TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

func FindStockWithMaximumIncrease(summaries []StockSummary) StockSummary {
	maxChange := 0.0
	maxSummary := StockSummary{}

	for _, summary := range summaries {
		change := summary.GetPercentChange()
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
		stocksByDay[util.TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

func (collection SummaryData) ForSymbolOnDay(symbol string, date time.Time) StockSummary {
	return collection.SummariesBySymbolByDay[symbol][date]
}

func (collection SummaryData) BeforeDay(day time.Time) []StockSummary {
	var stocks []StockSummary
	for currentDay, stocksOnDay := range collection.SummariesByDay {
		if currentDay.Before(day) {
			stocks = append(stocks, stocksOnDay...)
		}
	}

	return stocks
}

func (collection SummaryData) InDateRange(start time.Time, end time.Time) []StockSummary {
	var stocks []StockSummary
	for start.Before(end) {
			stocks = append(stocks, GetDailyStockSummaryData().OnDay(start)...)
			start = start.Add(24 * time.Hour)
	}

	return stocks
}

func (collection SummaryData) ForSymbolOnOtherDay(symbol string, date time.Time, days int64) StockSummary {
	duration := time.Duration(24 * days) * time.Hour
	nextDate := date.Add(duration)
	tries := 0
	for tries < 10 {
		nextDate = nextDate.Add(24 * time.Hour)
		tries++

		if collection.SummariesBySymbolByDay[symbol][nextDate].Symbol != "" {
			return collection.SummariesBySymbolByDay[symbol][nextDate]
		}
	}
	return collection.SummariesBySymbolByDay[symbol][date]
}

func (collection SummaryData) ForSymbolOnOtherDays(symbol string, date time.Time, days int) []StockSummary {
	tries := 0
	var summaries []StockSummary
	absDays := days
	if absDays < 0 {
		absDays *= -1
	}

	for tries < 10 && len(summaries) < absDays {
		if days > 0 {
			date = date.Add(24 * time.Hour)
		} else {
			date = date.Add(-24 * time.Hour)
		}

		if collection.SummariesBySymbolByDay[symbol][date].Symbol != "" {
			summaries = append(summaries, collection.SummariesBySymbolByDay[symbol][date])
		} else {
			tries++
		}
	}
	return summaries
}

func (collection SummaryData) OnDay(date time.Time) []StockSummary {
	return collection.SummariesByDay[date]
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
	SummariesByDay         map[time.Time][]StockSummary
	SummariesBySymbol      map[string][]StockSummary
	SummariesBySymbolByDay map[string]map[time.Time]StockSummary
}

type DailyStocksBySymbol map[string]map[string]StockSummary
