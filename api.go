package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
  "errors"
  "strings"
)

type ResponseMetaData struct {
	Information   string `json:"1. Information"`
	Ticker        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. US/Eastern"`
}

type RawWeeklyQuote struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	Adjusted         string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	Dividend         string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient"`
}

type WeeklyQuoteResponse struct {
	MetaData        ResponseMetaData         `json:"Meta Data"`
	TimeSeriesWeekly map[string]RawWeeklyQuote `json:"Weekly Adjusted Time Series"`
}

func request(f string, symbol string, outputsize string, tries int) ([]byte, error) {
	baseURL := "https://www.alphavantage.co"

	url := fmt.Sprintf(""+
		"%s/query?function=%s&symbol=%s&outputsize=%s&apikey=%s",
		baseURL,
		f,
		symbol,
		outputsize,
		API_KEY)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		panic("request failed for " + f + " " + symbol)
	}

	return body, nil
}

func TransformRawStock(symbol string, date string, rawStock RawWeeklyQuote) Quote {
	stock := Quote{}
	parsedDate, _ := time.Parse("2006-01-02", date)
	open, _ := strconv.ParseFloat(rawStock.Open, 64)
	high, _ := strconv.ParseFloat(rawStock.High, 64)
	low, _ := strconv.ParseFloat(rawStock.Low, 64)
	closePrice, _ := strconv.ParseFloat(rawStock.Close, 64)
	adjusted, _ := strconv.ParseFloat(rawStock.Adjusted, 64)
	volume, _ := strconv.ParseFloat(rawStock.Volume, 64)
	dividend, _ := strconv.ParseFloat(rawStock.Dividend, 64)

	stock.Ticker = symbol
	stock.Date = parsedDate
	stock.High = high
	stock.Low = low
	stock.Open = open
	stock.Close = closePrice
	stock.AdjustedClose = adjusted
	stock.Volume = volume
	stock.Dividend = dividend

	return stock
}

func GetWeeklyQuotes(symbol string) ([]Quote, error) {
	body, err := request("TIME_SERIES_WEEKLY_ADJUSTED", symbol, "full", 10)

	if err != nil {
		return nil, err
	}

	var response = new(WeeklyQuoteResponse)
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	var stocks []Quote

	for date, rawStock := range response.TimeSeriesWeekly {
		stock := TransformRawStock(response.MetaData.Ticker, date, rawStock)
		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
    
    if strings.Contains(string(body), "Invalid API call") {
      return nil, errors.New("invalid API call, ticker may not exist: " + symbol)
    }

		fmt.Println("waiting a bit: got no data for " + symbol)
		time.Sleep(60 * time.Second)
		fmt.Println("trying again: " + symbol)
		return GetWeeklyQuotes(symbol)
	}

	return stocks, nil
}


