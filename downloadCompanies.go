package main

import (
  "database/sql"
  "fmt"
	"encoding/csv"
	"io"
  "io/ioutil"
	"strconv"
  "net/http"
  "bytes"
)


type Company struct {
	Symbol       string
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

func ParseCompaniesFromCSV(ioReader io.Reader) map[string]Company {
	companies := make(map[string]Company)

	reader := csv.NewReader(ioReader)
	reader.Comma = ','

	// Skip headers
	_, err := reader.Read()
	if err != nil {
		panic(err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		symbol := record[0]
		name := record[1]
		lastTrade, _ := strconv.ParseFloat(record[2], 64)
		marketCap, _ := strconv.ParseFloat(record[3], 64)
		adr := record[4]
		tso := record[5]
		ipoyear, _ := strconv.ParseInt(record[6], 10, 64)
		sector := record[6]
		industry := record[7]
		summaryQuote := record[8]

		company := Company{
			symbol,
			name,
			lastTrade,
			marketCap,
			adr, tso,
			ipoyear,
			sector,
			industry,
			summaryQuote}

		companies[symbol] = company
	}

	return companies
}

func insertCompany(db *sql.DB, company Company) {
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
    company.Symbol,
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

func main() {
  url := "http://www.nasdaq.com/screening/companies-by-industry.aspx?render=download"
  resp, err := http.Get(url)

  if err != nil {
    panic("Could not download spreadsheet")
  }

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    panic("Could not parse body")
  }

  reader := bytes.NewReader(body)
  companies := ParseCompaniesFromCSV(reader)

  db := GetHandle()

  for ticker, company := range companies {
    fmt.Printf("Inserting %s\n", ticker)
    insertCompany(db, company)
  }

  // db := GetHandle()
  // fmt.Printf("%v", rows[0])
}

