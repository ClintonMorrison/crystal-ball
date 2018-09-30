package main

type QuoteClassifier func(*Quote) string

func GradeQuoteClassifier(q *Quote) string {
	return q.GetGrade()
}


type NGramModel struct {
	Models map[string]*NGramUniverse
	LabelForQuote QuoteClassifier
	Parameter string
	MaxN int
}

func CreateNGramModel(labelForQuote QuoteClassifier, parameter string, maxN int) *NGramModel {
	m := NGramModel{}
	m.Models = make(map[string]*NGramUniverse, 0)
	m.LabelForQuote = labelForQuote
	m.Parameter = parameter
	m.MaxN = maxN
	return &m
}

func (m *NGramModel) ParameterHashForCompany(company *Company) string {
	if m.Parameter == "INDUSTRY" {
		return company.Industry
	} else if m.Parameter == "SECTOR" {
		return company.Sector
	} else if m.Parameter == "TICKER" {
		return company.Ticker
	}

	return ""
}


func (m *NGramModel) GetModel(company *Company) *NGramUniverse {
	hash := m.ParameterHashForCompany(company)

	if m.Models[hash] == nil {
		universe := CreateUniverse(m.MaxN)
		m.Models[hash] = &universe
	}

	return m.Models[hash]
}


func (m *NGramModel) AddCase(company *Company, document string) {
	model := m.GetModel(company)
	model.AddString(document)
}

func (m *NGramModel) PredictNext(company *Company, document string) (string, float64) {
	model := m.GetModel(company)

	ngramStart := len(document) - m.MaxN + 1
	ngram := document[ngramStart:]

	return model.ProbableCompletion(ngram)
}