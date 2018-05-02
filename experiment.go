package main

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

type Strategy func(state *ExperimentState) []Order

type Experiment struct {
	params  ExperimentParams
	state   ExperimentState
	strategy Strategy
}

func initialStateFromParams(params ExperimentParams) *ExperimentState {
	state := ExperimentState{}
	state.Balance = params.InitialBalance
	state.Day = *&params.StartDay
	state.Portfolio = make(Portfolio)
	state.Params = params
	return &state
}

func CreateExperiment(params ExperimentParams, strategy Strategy) *Experiment {
	return &Experiment{
		params,
		*initialStateFromParams(params),
		strategy,
	}
}

func (experiment Experiment) Run() {
	for experiment.state.Day.Before(experiment.params.EndDay) {
		fmt.Println("")
		fmt.Println("", experiment.state.Day.Format("2006-01-02"))
		orders := experiment.strategy(&experiment.state)
		for _, order := range orders {
			experiment.state.applyOrder(order)
		}

		experiment.state.reportCurrentState()

		experiment.state.Day = experiment.state.Day.Add(time.Hour * 24)
	}

	experiment.state.reportSummary()
	fmt.Println("\n\nDONE")
}
