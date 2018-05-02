package trading

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

type StockInDB struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Symbol           string
	Date             time.Time
	Open             float64
	High             float64
	Low              float64
	Close            float64
	Adjusted         float64
	Volume           float64
	DividendAmount   float64
	SplitCoefficient float64
}

func connect() *mgo.Session {
	session, err := mgo.Dial(MongoURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func getStockyDailyCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("trading").C("stocksDaily")
}

func getCompaniesCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("trading").C("companies")
}

func GetCompanies() []Company {
	session := connect()
	defer session.Close()
	c := getCompaniesCollection(session)

	var results []Company
	c.Find(nil).All(&results)
	return results
}

func GetCompaniesBySmybol() map[string]Company {
	companies := GetCompanies()
	companiesBySymbol := make(map[string]Company, len(companies))
	for _, company := range companies {
		companiesBySymbol[company.Symbol] = company
	}

	return companiesBySymbol
}

func GetStocksDailyData() []Stock {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []Stock
	c.Find(nil).All(&results)
	return results
}

var cachedStocksForDay map[string][]Stock

func GetAvailableStocksForDay(date time.Time) []Stock {
	if cachedStocksForDay == nil {
		cachedStocksForDay = make(map[string][]Stock)
	}
	cachedResult, cacheHit := cachedStocksForDay[TimeToString(date)]
	if cacheHit {
		fmt.Println("****** CACHE HIT")
		return cachedResult
	}


	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []Stock
	c.Find(bson.M{"date": date}).Sort("symbol").All(&results)
	cachedStocksForDay[TimeToString(date)] = results
	return results
}

func GetAvailableStocksBySymbol(date time.Time) map[string]Stock {
	stocks := GetAvailableStocksForDay(date)
	stocksBySmybol := make(map[string]Stock)
	for _, stock := range stocks {
		stocksBySmybol[stock.Symbol] = stock
	}

	return stocksBySmybol
}

func GetPricesForStock(symbol string) []Stock {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []Stock
	c.Find(bson.M{"symbol": symbol}).Sort("date").All(&results)
	return results
}

var cachedStockForDay map[string]Stock

func GetStockForDay(symbol string, time time.Time) Stock {
	if cachedStockForDay == nil {
		cachedStockForDay = make(map[string]Stock)
	}
	cachedResult, cacheHit := cachedStockForDay[symbol + TimeToString(time)]
	if cacheHit {
		fmt.Println("****** CACHE HIT")
		return cachedResult
	}

	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var result Stock
	err := c.Find(bson.M{"symbol": symbol, "date": time}).One(&result)
	if err != nil {
		fmt.Println(err)
		return Stock{}
	}
	cachedStockForDay[symbol + TimeToString(time)] = result
	return result
}

func GetPricesForStockByDay(symbol string) map[string]Stock {
	stocks := GetPricesForStock(symbol)
	stocksByDay := make(map[string]Stock, len(stocks))
	for _, stock := range stocks {
		stocksByDay[TimeToString(stock.Date)] = stock
	}

	return stocksByDay
}

func AddStockDay(s Stock) {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	c.RemoveAll(bson.M{"symbol": s.Symbol, "date": s.Date})

	err := c.Insert(&s)
	if err != nil {
		panic(err)
	}
}

func BatchAddStockDay(stocks []Stock) {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	bulk := c.Bulk()

	/*
	for _, s := range stocks {
		bulk.RemoveAll(bson.M{"symbol": s.Symbol, "date": s.Date})
	}
	bulk.Run()
	*/

	bulk = c.Bulk()
	for _, s := range stocks {
		bulk.Insert(s)
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
		fmt.Println("Inserting company", s.Symbol)
		bulk.Insert(s)
	}

	_, err := bulk.Run()
	if err != nil {
		panic(err)
	}
}
