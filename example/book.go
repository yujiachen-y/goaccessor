package main

//go:generate goaccessor -t bookStoreName -g
const bookStoreName = "Go's Book Store"

//go:generate goaccessor --target Book --getter --setter
type Book struct {
	Title  string
	Author string
}

//go:generate goaccessor --target books --getter --setter
var books = map[string]*Book{
	"book 1": {Title: "title 1", Author: "author 1"},
	"book 2": {Title: "title 2", Author: "author 2"},
	"book 3": {Title: "title 3", Author: "author 3"},
}

//go:generate goaccessor --target bestSellingBook --field --getter --include Author --prefix BestSelling
var bestSellingBook = &Book{
	Title:  "Best Selling Title",
	Author: "Best Selling Author",
}

func main() {}
