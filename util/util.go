package util

import "time"

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDay(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func Sum(nums []float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}

	return total
}

func Avg(nums []float64) float64 {
	return Sum(nums) / float64(len(nums))
}

func CountAboveThreshold(threshold float64, nums []float64) int {
	count := 0
	for _, n := range nums {
		if n > threshold {
			count++
		}
	}

	return count
}