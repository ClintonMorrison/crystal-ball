package main

import "fmt"

type Order struct {
	Action
	Symbol   string
	Quantity float64
}

func (order Order) reportCost(cost float64) {
	var action string
	switch order.Action {
	case BUY:
		action = "BUY"
	case SELL:
		action = "SELL"
	}

	fmt.Printf("   --> %s x %.2f shares of %s for $%.2f\n", action, order.Quantity, order.Symbol, cost)
}
