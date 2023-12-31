package models

type Data struct {
	Quotes       []Quote
	Quote        string
	Author       string
	IsSubscribed bool
	Subscribers  []EmailData
	CSRFToken    string
	Bool         bool
	IntMap       map[string]int
	StringMap    map[string]string
	Error        error
	IsError      bool
}

type Quote struct {
	ID     int
	Quote  string
	Author string
}

type EmailData struct {
	Email string
	Name  string
}

var Subscriber EmailData

type WelcomeMail struct {
	To   string
	Name string
}
