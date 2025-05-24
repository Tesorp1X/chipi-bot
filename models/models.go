package models

type Check struct {
	Name  string
	Owner string
}

type Item struct {
	CheckId int64
	Name    string
	Owner   string
	Price   float64
}
