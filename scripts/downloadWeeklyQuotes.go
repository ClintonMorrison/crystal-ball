package scripts

import (
  "fmt"
)

func main() {
  db := GetHandle()
  tickers := GetAllCompanyTickers(db)
  start := false
  lastSymbol := "AKO"
  for i, ticker := range tickers {
    fmt.Printf("[%4d/%4d] - %s\n", i, len(tickers), ticker)

    if ticker == lastSymbol {
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

