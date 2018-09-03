package data

import (
	"encoding/csv"
	"github.com/globalsign/mgo/bson"
	"io"
	"os"
	"strconv"
)

type Company struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Symbol       string
	Name         string
	LastSale     float64
	MarketCap    float64
	ADR          string
	TSO          string
	IPOyear      int64
	Sector       string
	Industry     string
	SummaryQuote string
}

func ParseCompaniesFromCSV() map[string]Company {
	companies := make(map[string]Company)

	file, err := os.Open("data/companylist.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','

	// Skip headers
	_, err = reader.Read()
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
			"",
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

func GetSymbols(companies map[string]Company) []string {
	keys := make([]string, 0, len(companies))
	for k, _ := range companies {
		keys = append(keys, k)
	}

	return keys
}

func GetCompaniesBySmybol() map[string]Company {
	companies := GetCompanies()
	companiesBySymbol := make(map[string]Company, len(companies))
	for _, company := range companies {
		companiesBySymbol[company.Symbol] = company
	}

	return companiesBySymbol
}
