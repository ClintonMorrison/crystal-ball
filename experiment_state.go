package main

import (
	"fmt"
	"time"
)

type Portfolio map[string]float64 // symbol => quantity held

type ExperimentState struct {
	Balance float64
	Day     time.Time
	Portfolio
	Params ExperimentParams
}

func (state ExperimentState) lookupPrice(symbol string, date time.Time) float64 {
	price := GetDailyStockSummaryData().forSymbolOnDay(symbol, date).Close // GetDailySummaryForStock(symbol, date).Close // state.Params.DailyStocksBySymbol[symbol][dateString].Close

	tries := 0

	for price == 0 && tries < 10 {
		tries++
		price = GetDailySummaryForStock(symbol, date.AddDate(0, 0, -1*tries)).Close
	}

	return price
}

func (state ExperimentState) reportCurrentState() {
	fmt.Printf("  balance: $%.2f\n", state.Balance)
	for symbol, qty := range state.Portfolio {
		if qty > 0 {
			fmt.Printf("  %7s: %.2f\n", symbol, qty)
		}
	}
}

func (state ExperimentState) getPortfolioValue() float64 {
	total := 0.0
	for symbol, qty := range state.Portfolio {
		total += qty * state.lookupPrice(symbol, state.Day)
	}

	return total
}

func (state ExperimentState) getTotalValue() float64 {
	return state.Balance + state.getPortfolioValue()
}

func (state ExperimentState) reportSummary() {
	fmt.Println("\n-----------\n\nSUMMARY")
	fmt.Printf("  balance: $%.2f\n", state.Balance)
	for symbol, qty := range state.Portfolio {
		fmt.Printf("  %7s: %.2f\n", symbol, qty)
	}

	initialValue := state.Params.InitialBalance
	finalValue := state.getTotalValue()
	change := finalValue - initialValue

	fmt.Println()
	fmt.Printf("  Initial: $%.2f\n", initialValue)
	fmt.Printf("    Final: $%.2f\n", finalValue)
	fmt.Printf("   Change: $%.2f\n", change)
	fmt.Printf("\n   Profit: %.2f%%\n", change/initialValue*100.0)
}

func (state ExperimentState) applyOrder(order Order) ExperimentState {
	pricePerShare := state.lookupPrice(order.Symbol, state.Day)
	if pricePerShare == 0 {
		fmt.Println("WARN: Skipping trade. Price is 0 for " + order.Symbol + ", " + TimeToString(state.Day))

		// Stock no longer exists? do nothing
		return state
	}
	multiplier := 1.0

	if order.Action == SELL {
		multiplier = -1.0
	}

	newQuantity := state.Portfolio[order.Symbol] + (order.Quantity * multiplier)
	cost := order.Quantity * pricePerShare
	newBalance := state.Balance - (cost * multiplier)
	if newQuantity >= 0 && newBalance >= 0 {
		order.reportCost(cost)
		state.Portfolio[order.Symbol] = newQuantity
		if newQuantity == 0 {
			delete(state.Portfolio, order.Symbol)
		}

		if state.Params.TransactionFee > 0 {
			fmt.Println("applying transaction fee of", state.Params.TransactionFee)
		}
		state.Balance = newBalance - state.Params.TransactionFee
	}

	return state
}
