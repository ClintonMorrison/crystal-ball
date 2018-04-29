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
	session, err := mgo.Dial("localhost:27017")
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

func GetStocksDailyData() []Stock {
	session := connect()
	defer session.Close()
	c := getStockyDailyCollection(session)

	var results []Stock
	c.Find(nil).All(&results)
	return results
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

	for _, s := range stocks {
		bulk.RemoveAll(bson.M{"symbol": s.Symbol, "date": s.Date})
	}
	bulk.Run()

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
