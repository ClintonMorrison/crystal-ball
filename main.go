package main

import (
	"fmt"
	"time"
	"trading/trading"
)

func parseDay(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func doNothing(state trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	fmt.Println("running")
	return orders
}

func buyAll(state *trading.ExperimentState) []trading.Order {
	var orders []trading.Order
	orders = append(orders, trading.Order{trading.BUY, "MSFT", 1.0})
	return orders
}

/*
func getAvailableSymbols(state *ExperimentState) []string {

}
*/

func main() {

	/*
		companies := GetCompanies()
		symbols := GetSymbols(companies)
		fmt.Println(symbols)

		dailyStocksBySymbol := GetPricesByDayForStocks([]string{
			symbols[0], symbols[10], symbols[100], symbols[105], symbols[200], symbols[500], symbols[1000]})


		params := ExperimentParams{}
		params.DailyStocksBySymbol = dailyStocksBySymbol
		params.InitialBalance = 1000
		params.StartDay = parseDay("2015-01-01")
		params.EndDay = parseDay("2018-01-28")
		fmt.Printf("%#v\n\n", dailyStocksBySymbol)

		RunExperiment(params, buyAll)
	*/

	trading.AddStock(trading.Stock{Symbol: "TEST"})
	fmt.Println(trading.GetDBData())
}
