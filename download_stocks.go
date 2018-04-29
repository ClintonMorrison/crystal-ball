package main

import (
	"fmt"
	"stock-analysis/trading"
)

func SaveStockData(symbol string) {
	fmt.Println("Fetching data for " + symbol)
	stockDays, err := trading.GetDailyStockData(symbol)
	if err != nil {
		panic(err)
	}

	trading.BatchAddStockDay(stockDays)
}

func SaveCompanies(companiesBySymbol map[string]trading.Company) {
	fmt.Println("Saving companies...")
	var companies []trading.Company
	for _, company := range companiesBySymbol {
		companies = append(companies, company)
	}

	trading.BatchAddCompany(companies)
}

func main() {
	companiesBySymbol := trading.GetCompanies()
	SaveCompanies(companiesBySymbol)

	for symbol, _ := range companiesBySymbol {
		SaveStockData(symbol)
	}
}
