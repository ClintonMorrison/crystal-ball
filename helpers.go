package main

import (
  "fmt"
  "time"
)

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDays(t time.Time, days float64) time.Time {
  return t.Add(time.Duration(days * 24) * time.Hour)
}

func AddWeek(t time.Time) time.Time {
  return AddDays(t, 7)
}


func ParseDateString(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func ListMin(values []float64) float64 {
  min := values[0]
  for _, val := range values {
    if val < min {
      min = val
    }
  }

  return min
}


func ListMax(values []float64) float64 {
  max := values[0]
  for _, val := range values {
    if val > max {
      max = val
    }
  }

  return max
}

func ListCountInRange(values [] float64, min float64, max float64) int {
  count := 0
  for _, val := range values {
    if val >= min && val <= max {
      count++
    }
  }

  return count
}


func PrintSequenceStats(values []float64) {
  min := ListMin(values)
  max := ListMax(values)

  fmt.Printf("Min: %.2f\n", min)
  fmt.Printf("Max: %.2f\n", max)

  bucketSize := 0.05
  start := 1.0
  end := -1.0
  current := start
  
  for current > end {
    count := ListCountInRange(values, current - bucketSize, current)
    fmt.Printf("[%2.2f - %2.2f] %d\n", current, current - bucketSize, count)
    current -= bucketSize
  }
}

func GetGradeString(quotes []*Quote) string {
  gradeString := ""

  for _, quote := range quotes {
    grade := quote.GetGrade()
    gradeString = gradeString + grade
  }

  return gradeString
}