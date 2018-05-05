package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"stock-analysis/config"
)

// Idea
// Create "strategies" and then apply them on historical data
// report gains and percentage of "windows" where it works well
// report most catastrpohic losses of windows too
//
// e.g.

type ResponseMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. US/Eastern"`
}

type RawDailyStock struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	Adjusted         string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient"`
}

type DailyStockResponse struct {
	MetaData        ResponseMetaData         `json:"Meta Data"`
	TimeSeriesDaily map[string]RawDailyStock `json:"Time Series (Daily)"`
}

func request(f string, symbol string, outputsize string) ([]byte, error) {
	baseURL := "https://www.alphavantage.co"

	url := fmt.Sprintf(""+
		"%s/query?function=%s&symbol=%s&outputsize=%s&apikey=%s",
		baseURL,
		f,
		symbol,
		outputsize,
		config.ApiKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func TransformRawStock(symbol string, date string, rawStock RawDailyStock) StockSummary {
	stock := StockSummary{}
	parsedDate, _ := time.Parse("2006-01-02", date)
	open, _ := strconv.ParseFloat(rawStock.Open, 64)
	high, _ := strconv.ParseFloat(rawStock.High, 64)
	low, _ := strconv.ParseFloat(rawStock.Low, 64)
	closePrice, _ := strconv.ParseFloat(rawStock.Close, 64)
	adjusted, _ := strconv.ParseFloat(rawStock.Adjusted, 64)
	volume, _ := strconv.ParseFloat(rawStock.Volume, 64)
	dividendAmount, _ := strconv.ParseFloat(rawStock.DividendAmount, 64)
	splitCoefficient, _ := strconv.ParseFloat(rawStock.SplitCoefficient, 64)

	stock.Symbol = symbol
	stock.Date = parsedDate
	stock.High = high
	stock.Low = low
	stock.Open = open
	stock.Close = closePrice
	stock.AdjustedClose = adjusted
	stock.Volume = volume
	stock.DividendAmount = dividendAmount
	stock.SplitCoefficient = splitCoefficient

	return stock
}

func GetDailyStockData(symbol string) ([]StockSummary, error) {
	body, err := request("TIME_SERIES_DAILY_ADJUSTED", symbol, "full") // "full" "compact"

	if err != nil {
		return nil, err
	}

	var response = new(DailyStockResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var stocks []StockSummary

	for date, rawStock := range response.TimeSeriesDaily {
		stock := TransformRawStock(response.MetaData.Symbol, date, rawStock)
		stocks = append(stocks, stock)
	}

	// fmt.Printf("%#v", response)
	return stocks, nil
}

func GetPricesByDayForStock(symbol string) map[string]StockSummary {
	stocks, err := GetDailyStockData(symbol)
	if err != nil {
		panic(err)
	}

	return GroupStocksByDay(stocks)
}

func GetPricesByDayForStocks(symbols []string) DailyStocksBySymbol {
	dailyStocksBySymbol := make(DailyStocksBySymbol)

	for _, symbol := range symbols {
		dailyStocksBySymbol[symbol] = GetPricesByDayForStock(symbol)
	}

	return dailyStocksBySymbol
}
