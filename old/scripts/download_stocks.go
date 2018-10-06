package main

import (
	"fmt"
	"stockAnalysis/data"
	// "time"
)

func SaveStockData(symbol string) {
	stockDays, err := data.GetDailyStockData(symbol)
	// time.Sleep(20 * time.Second) // Max of 5 requests per minute

	if err != nil {
		panic(err)
	}

	data.BatchAddStockDailySummary(symbol, stockDays)
}

func SaveCompanies(companiesBySymbol map[string]data.Company) {
	fmt.Println("Saving companies...")
	var companies []data.Company
	for _, company := range companiesBySymbol {
		companies = append(companies, company)
	}

	data.BatchAddCompany(companies)
}

func main() {

	companiesBySymbol := data.ParseCompaniesFromCSV()
	SaveCompanies(companiesBySymbol)

	index := 0
	for symbol, _ := range companiesBySymbol {
		index += 1
		fmt.Printf("Fetching data for %s [%d of %d] --> ", symbol, index, len(companiesBySymbol))
		SaveStockData(symbol)
	}
}
