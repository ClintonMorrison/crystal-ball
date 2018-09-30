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

func (q *Quote) Print() {
	fmt.Println(q.Ticker)
	fmt.Println("-------")
	fmt.Printf("Date: %s\n", TimeToString(q.Date))
	fmt.Printf("Open: $%.2f\n", q.Open)
	fmt.Printf("Close: $%.2f\n", q.Close)
	fmt.Printf("Low: $%.2f\n", q.Low)
	fmt.Printf("High: $%.2f\n", q.High)
	fmt.Printf("Volume: %.2f\n", q.Volume)
	fmt.Printf("Change: %%%.2f\n", q.GetPercentChange() * 100.0)
	fmt.Printf("Violatility: %%%.2f\n", q.GetPercentVolatility() * 100.0)
}

func (q *Quote) GetPercentChange() float64 {
  return (q.Close - q.Open) / q.Open
}

func (q *Quote) GetPercentVolatility() float64 {
  return (q.High - q.Low) / q.Open
}

func (q *Quote) GetGrade() string {
  change := q.GetPercentChange()

  if change > 0.01 {
  	return "U"
	} else if change < -0.01 {
		return "D"
	} else {
		return "_"
	}



  if change >= 0.05 {
    return "A"
  } else if change >= 0.025 {
    return "B"
  } else if change >= 0 {
    return "C"
  } else if change >= -0.025 {
    return "D"
  } else if change >= -0.05 {
    return "E"
  } else {
    return "F"
  }
}

func rowsToQuotes(rows *sql.Rows) []*Quote {
	quotes := make([]*Quote, 0)
	for rows.Next() {
		quote := Quote{}
		err := rows.Scan(
			&quote.Ticker,
			&quote.Date,
			&quote.Open,
			&quote.High,
			&quote.Low,
			&quote.Close,
			&quote.AdjustedClose,
			&quote.Volume,
			&quote.Dividend)

		if err != nil {
			fmt.Println(err)
			continue
		}
		quotes = append(quotes, &quote)
	}

	return quotes
}

func rowToQuote(row *sql.Row) (*Quote, error) {
  quote := Quote{}

  err := row.Scan(
    &quote.Ticker,
    &quote.Date,
    &quote.Open,
    &quote.High,
    &quote.Low,
    &quote.Close,
    &quote.AdjustedClose,
    &quote.Volume,
    &quote.Dividend)

  if err != nil {
    return nil, err
  }

  return &quote, nil;
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

func GetQuoteForTickerOnWeek(db *sql.DB, ticker string, date time.Time) (*Quote, error) {
  weekStartDate := AddDays(date, -3.5)
  weekEndDate := AddDays(date, 3.5)

  // fmt.Printf("WHERE ticker = %s. date range %s < ... < %s\n", ticker,
  // TimeToString(weekStartDate), TimeToString(weekEndDate))
  row := db.QueryRow(`
    SELECT
      ticker,
      date,
      open,
      high,
      low,
      close,
      adjusted_close,
      volume,
      dividend FROM
    quotes WHERE
      ticker = $1 AND
      date >= $2 AND
      date <= $3`,
    ticker,
    weekStartDate,
    weekEndDate)

  return rowToQuote(row)
}


func GetAllQuotesForTicker(db *sql.DB, ticker string) ([]*Quote) {
  rows, err := db.Query(`
    SELECT
      ticker,
      date,
      open,
      high,
      low,
      close,
      adjusted_close,
      volume,
      dividend FROM
    quotes WHERE
      ticker = $1
    ORDER BY date ASC`,
    ticker)

  if err != nil {
  	panic(err)
	}

	defer rows.Close()

	return rowsToQuotes(rows)
}


