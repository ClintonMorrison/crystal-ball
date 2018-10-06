package main

import (
  "fmt"
	"encoding/csv"
	"io"
  "io/ioutil"
	"strconv"
  "net/http"
  "bytes"
  "strings"
)


func parseCompaniesFromCSV(ioReader io.Reader) map[string]Company {
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
		symbol := strings.TrimSpace(record[0])
    symbol = strings.Split(symbol, "^")[0]
    symbol = strings.Split(symbol, ".")[0]
		name := strings.TrimSpace(record[1])
		lastTrade, _ := strconv.ParseFloat(record[2], 64)
		marketCap, _ := strconv.ParseFloat(record[3], 64)
		adr := strings.TrimSpace(record[4])
		tso := strings.TrimSpace(record[5])
		ipoyear, _ := strconv.ParseInt(record[6], 10, 64)
		sector := strings.TrimSpace(record[6])
		industry := strings.TrimSpace(record[7])
		summaryQuote := strings.TrimSpace(record[8])

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

func DownloadCompanies() {
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
  companies := parseCompaniesFromCSV(reader)

  db := GetHandle()

  for ticker, company := range companies {
    fmt.Printf("Inserting %s\n", ticker)
    InsertCompany(db, company)
  }
}

