package trading

import (
	"fmt"
	"time"
)

type ExperimentParams struct {
	InitialBalance      float64
	StartDay            time.Time
	EndDay              time.Time
	CompaniesBySymbol	  map[string]Company
}

type Action int

const (
	BUY Action = iota
	SELL
)

type Order struct {
	Action
	Symbol   string
	Quantity float64
}

type Portfolio map[string]float64 // symbol => quantity held

type ExperimentState struct {
	Balance float64
	Day     time.Time
	Portfolio
	Params ExperimentParams
}

type Stategy func(state *ExperimentState) []Order

func initialStateFromParams(params ExperimentParams) *ExperimentState {
	state := ExperimentState{}
	state.Balance = params.InitialBalance
	state.Day = *&params.StartDay
	state.Portfolio = make(Portfolio)
	state.Params = params
	return &state
}

func lookupPrice(state *ExperimentState, symbol string, date time.Time) float64 {
	price := GetStockForDay(symbol, date).Close // state.Params.DailyStocksBySymbol[symbol][dateString].Close

	tries := 0

	for price == 0 && tries < 10 {
		tries++
		price = GetStockForDay(symbol, date.AddDate(0, 0, -1 * tries)).Close
	}

	return price
}

func reportState(state *ExperimentState) {
	fmt.Printf("  balance: $%.2f\n", state.Balance)
	for symbol, qty := range state.Portfolio {
		if qty > 0 {
			fmt.Printf("  %7s: %.2f\n", symbol, qty)
		}
	}
}

func getPortfolioValue(state *ExperimentState) float64 {
	total := 0.0
	for symbol, qty := range state.Portfolio {
		total += qty * lookupPrice(state, symbol, state.Day)
	}

	return total
}

func getTotalValue(state *ExperimentState) float64 {
	return state.Balance + getPortfolioValue(state)
}

func reportSummary(state *ExperimentState) {
	fmt.Println("\n-----------\n\nSUMMARY")
	fmt.Printf("  balance: $%.2f\n", state.Balance)
	for symbol, qty := range state.Portfolio {
		fmt.Printf("  %7s: %.2f\n", symbol, qty)
	}

	initialValue := state.Params.InitialBalance
	finalValue := getTotalValue(state)
	change := finalValue - initialValue

	fmt.Println()
	fmt.Printf("  Initial: $%.2f\n", initialValue)
	fmt.Printf("    Final: $%.2f\n", finalValue)
	fmt.Printf("   Change: $%.2f\n", change)
	fmt.Printf("\n   Profit: %.2f%%\n", change/initialValue*100.0)

}

func reportOrder(order Order, cost float64) {
	var action string
	switch order.Action {
	case BUY:
		action = "BUY"
	case SELL:
		action = "SELL"
	}

	fmt.Printf("   --> %s x %.2f shares of %s for $%.2f\n", action, order.Quantity, order.Symbol, cost)

}

func applyOrder(state *ExperimentState, order Order) {
	pricePerShare := lookupPrice(state, order.Symbol, state.Day)
	if pricePerShare == 0 {
		panic("Price is 0 for " + order.Symbol + ", " + TimeToString(state.Day))
	}
	multiplier := 1.0

	if order.Action == SELL {
		multiplier = -1.0
	}

	newQuantity := state.Portfolio[order.Symbol] + (order.Quantity * multiplier)
	cost := order.Quantity * pricePerShare
	newBalance := state.Balance - (cost * multiplier)
	if newQuantity >= 0 && newBalance >= 0 {
		reportOrder(order, cost)
		state.Portfolio[order.Symbol] = newQuantity
		state.Balance = newBalance
	}
}

func RunExperiment(params ExperimentParams, stategy Stategy) {
	state := initialStateFromParams(params)
	for state.Day.Before(params.EndDay) {
		fmt.Println("")
		fmt.Println("", state.Day.Format("2006-01-02"))
		orders := stategy(state)
		for _, order := range orders {
			applyOrder(state, order)
		}

		reportState(state)

		state.Day = state.Day.Add(time.Hour * 24)
	}

	reportSummary(state)
	fmt.Println("\n\nDONE")
}
