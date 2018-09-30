package main

import (
	"fmt"
	"sort"
)

type EvaluationCondition func(actual string, predicted string) bool

func isCorrect(actual string, predicted string) bool {
	return actual == predicted
}


func PrintEvaluationResult(field string, numerator int, denominator int) {
	percent := float64(numerator)/float64(denominator)*100.0
	fmt.Printf("%s: %2.2f%%  [%d/%d]\n", field, percent, numerator, denominator)
}

type EvaluationResults struct {
	actualLabels    []string
	predictedLabels []string
}

func (r *EvaluationResults) AddCase(actual string, predicted string) {
	r.actualLabels = append(r.actualLabels, actual)
	r.predictedLabels = append(r.predictedLabels, predicted)
}

func (r *EvaluationResults) init() {
	r.actualLabels = make([]string, 0)
	r.predictedLabels = make([]string, 0)
}

func (r *EvaluationResults) DistinctLabels() []string {
	labelMap := make(map[string]bool, 0)
	labels := make([]string, 0)

	for _, label := range r.actualLabels {
		labelMap[label] = true
	}

	for label := range labelMap {
		labels = append(labels, label)
	}

	sort.Strings(labels)

	return labels
}

func (r *EvaluationResults) CountCondition(condition EvaluationCondition) int {
	count := 0
	for i := range r.actualLabels {
		actual := r.actualLabels[i]
		predicted := r.predictedLabels[i]

		if condition(actual, predicted) {
			count++
		}
	}

	return count
}

func (r *EvaluationResults) CorrectlyPredictedForLabel(label string) (int, int) {
	wasPredictedAsLabel := func(actual string, predicted string) bool {
		return predicted == label
	}

	wasPredictedCorrectlyAsLabel := func(actual string, predicted string) bool {
		return predicted == label && isCorrect(actual, predicted)
	}

	totalCorrect := r.CountCondition(wasPredictedCorrectlyAsLabel)
	totalPredicted := r.CountCondition(wasPredictedAsLabel)

	return totalCorrect, totalPredicted
}

func (r *EvaluationResults) PercentCorrect() float64 {
	overallCorrect := r.CountCondition(isCorrect)
	overallTotal := len(r.predictedLabels)
	return float64(overallCorrect) / float64(overallTotal) * 100
}

func (r *EvaluationResults) PrintReport() {
	overallCorrect := r.CountCondition(isCorrect)
	overallTotal := len(r.predictedLabels)
	PrintEvaluationResult("Overall Correct", overallCorrect, overallTotal)

	labels := r.DistinctLabels()
	for _, label := range labels {
		correct, total := r.CorrectlyPredictedForLabel(label)
		PrintEvaluationResult(label, correct, total)
	}
}

