package main

import (
  "database/sql"
  "fmt"
)

type Company struct {
	Ticker       string
	Name         string
	LastSale     float64
	MarketCap    float64
	ADR          string
	TSO          string
	IPOYear      int64
	Sector       string
	Industry     string
	SummaryQuote string
}

func InsertCompany(db *sql.DB, company Company) {
  rows, err := db.Query(`
    INSERT INTO companies
    (
      ticker, name, industry, last_sale, market_cap, adr,
      tso, ipo_year, sector, summary_quote
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
      $9,
      $10
    )`,
    company.Ticker,
    company.Name,
    company.Industry,
    company.LastSale,
    company.MarketCap,
    company.ADR,
    company.TSO,
    company.IPOYear,
    company.Sector,
    company.SummaryQuote)

  if err != nil {
    fmt.Println(err)
  } else {
    defer rows.Close()
  }
}

func GetAllCompanyTickers(db *sql.DB) []string {
  tickers := make([]string, 0, 7000)

  rows, err := db.Query(`SELECT ticker FROM companies ORDER BY ticker ASC`)

  if err != nil {
    panic(err)
  }

  defer rows.Close()

  for rows.Next() {
    ticker := ""
    err := rows.Scan(&ticker)
    if err != nil {
      panic(err)
    }

    tickers = append(tickers, ticker)
  }

  return tickers
}

func GetCompanyByTicker(db *sql.DB, ticker string) *Company {
  company := Company{}

  row := db.QueryRow(`SELECT
    ticker,
    name,
    industry,
    last_sale,
    market_cap,
    adr,
    tso,
    ipo_year,
    sector,
    summary_quote
    FROM companies
    WHERE ticker = $1`,
    ticker)


  err := row.Scan(
    &company.Ticker,
    &company.Name,
    &company.Industry,
    &company.LastSale,
    &company.MarketCap,
    &company.ADR,
    &company.TSO,
    &company.IPOYear,
    &company.Sector,
    &company.SummaryQuote)

  if err != nil {
    panic(err)
  }

  return &company
}
