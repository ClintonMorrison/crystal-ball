package main

import (
  "fmt"
)

func DownloadWeeklyQuotes(startTicker string) {
  db := GetHandle()
  tickers := GetAllCompanyTickers(db)
  start := startTicker == ""

  for i, ticker := range tickers {
    fmt.Printf("[%4d/%4d] - %s\n", i, len(tickers), ticker)

    if !start && ticker == startTicker {
      start = true
    }

    if !start {
      continue
    }

    quotes, err := GetWeeklyQuotes(ticker)
    if err != nil {
      fmt.Println(err)
      continue
    }

    for _, quote := range quotes {
      InsertQuote(db, quote)
    }
  }

  fmt.Println("done")
}

