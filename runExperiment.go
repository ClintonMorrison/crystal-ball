package main

import (
  "fmt"
)

func main() {
  i := 1
  for i < 50 {
    runExperiment("ELLO", i)
    i++
  }
}

func runExperiment(ticker string, ngramSize int) {
  db := GetHandle()

  date := ParseDateString("2017-10-01")
  endDate := ParseDateString("2018-09-01")

	gradeString := ""
	correctUp := 0
	correctDown := 0
  totalUp := 0
  totalDown := 0

  for date.Before(endDate) {
    quote, err := GetQuoteForTickerOnWeek(db, ticker, date)
    date = AddWeek(date)

    if err != nil {
      fmt.Println(err)
      continue
    }

		actual := quote.GetGrade()

		if len(gradeString) > 0 {
			universe := CreateUniverse(ngramSize)
			universe.AddString(gradeString)
			prediction := universe.GenerateNextCharacter(gradeString)
			// fmt.Printf("[%s/%s] ", prediction, actual)

      actualWentDown := actual == "F" || actual == "E" || actual == "D"
      actualWentUp := actual == "A" || actual == "B" || actual == "C"
      if prediction == "F" && actualWentDown {
        correctDown += 1
      }

      if prediction == "A" && actualWentUp {
        correctUp += 1
      }

      if prediction == "A" {
        totalUp += 1
      }
      if prediction == "F" {
        totalDown += 1
      }
		}


		gradeString += actual

  }

  fmt.Printf("\n%s [n=%d] up correct: %2.2f / down correct: %2.2f \n",
    ticker,ngramSize, float64(correctUp)/float64(totalUp)*100.0,
    float64(correctDown)/float64(totalDown)*100.0)
}
