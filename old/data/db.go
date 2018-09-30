package data

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
	"stock-analysis/old/config"
)

func connect() *mgo.Session {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func getCollection(session *mgo.Session, collection string) *mgo.Collection {
	return session.DB("trading").C(collection)
}

// Companies
func getCompaniesCollection(session *mgo.Session) *mgo.Collection {
	return getCollection(session, "companies")
}

func GetCompanies() []Company {
	session := connect()
	defer session.Close()
	c := getCompaniesCollection(session)

	var results []Company
	c.Find(nil).All(&results)
	return results
}

// Stock prices
func getStockyDailyCollection(session *mgo.Session) *mgo.Collection {
	return getCollection(session, "stocksDaily")
}

func GetAllDailyStockPrices() []StockSummary {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []StockSummary
	c.Find(nil).All(&results)
	return results
}

func GetStockSummariesForDay(date time.Time) []StockSummary {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []StockSummary
	c.Find(bson.M{"date": date}).Sort("symbol").All(&results)
	return results
}

func GetDailySummariesForStock(symbol string) []StockSummary {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []StockSummary
	c.Find(bson.M{"symbol": symbol}).Sort("date").All(&results)
	return results
}

func GetDailySummaryForStock(symbol string, time time.Time) StockSummary {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var result StockSummary
	err := c.Find(bson.M{"symbol": symbol, "date": time}).One(&result)
	if err != nil {
		fmt.Println(err)
		return StockSummary{}
	}
	return result
}

func AddDailyStockSummary(s StockSummary) {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	c.RemoveAll(bson.M{"symbol": s.Symbol, "date": s.Date})

	err := c.Insert(&s)
	if err != nil {
		panic(err)
	}
}

func BatchAddStockDailySummary(symbol string, stocks []StockSummary) {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)
	c.RemoveAll(bson.M{"symbol": symbol})

	bulk := c.Bulk()
	for i, s := range stocks {
		bulk.Insert(s)

		if i % 100 == 0 {
			_, err := bulk.Run()
			if err != nil {
				panic(err)
			}
			bulk = c.Bulk()
		}

	}

	_, err := bulk.Run()
	if err != nil {
		panic(err)
	}
}

func BatchAddCompany(companies []Company) {
	session := connect()
	defer session.Close()
	c := getCompaniesCollection(session)

	bulk := c.Bulk()

	for _, company := range companies {
		bulk.RemoveAll(bson.M{"symbol": company.Symbol})
	}
	bulk.Run()

	bulk = c.Bulk()
	for _, s := range companies {
		bulk.Insert(s)
	}

	_, err := bulk.Run()
	if err != nil {
		panic(err)
	}
}
