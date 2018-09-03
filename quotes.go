package main

import (
  "database/sql"
	"time"
  "fmt"
)
type Quote struct {
	Ticker           string
	Date             time.Time
	Open             float64
	High             float64
	Low              float64
	Close            float64
	AdjustedClose    float64
	Volume           float64
	Dividend         float64
}

func InsertQuote(db *sql.DB, quote Quote) {
  rows, err := db.Query(`
    INSERT INTO quotes
    (
      ticker, volume, date, open, high,
      low, close, adjusted_close, dividend
    )
    
    VALUES (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8,
      $9
    )`,
    quote.Ticker,
    quote.Volume,
    quote.Date,
    quote.Open,
    quote.High,
    quote.Low,
    quote.Close,
    quote.AdjustedClose,
    quote.Dividend)

  if err != nil {
    fmt.Print(".")
  } else {
    defer rows.Close()
  }
}

