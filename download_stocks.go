package main

import (
	"fmt"
	"stock-analysis/trading"
)

func SaveStockData(symbol string) {
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
	companiesBySymbol := trading.ParseCompaniesFromCSV()
	SaveCompanies(companiesBySymbol)

	index := 0
	for symbol, _ := range companiesBySymbol {
		index += 1
		fmt.Printf("Fetching data for %s [%d of %d] --> ", symbol, index, len(companiesBySymbol))
		SaveStockData(symbol)
	}
}
