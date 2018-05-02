package main

import "time"

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDay(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
