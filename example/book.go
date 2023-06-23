package main

import "fmt"

//go:generate goaccessor -t bookStoreName -g
const bookStoreName = "Go's Book Store"

//go:generate goaccessor --target Book --getter --setter
type Book struct {
	title  string
	author string
}

//go:generate goaccessor --target books --getter --setter
var books = map[string]*Book{
	"book 1": {title: "title 1", author: "author 1"},
	"book 2": {title: "title 2", author: "author 2"},
	"book 3": {title: "title 3", author: "author 3"},
}

//go:generate goaccessor --target bestSellingBook --field --getter --include author --prefix BestSelling
var bestSellingBook = &Book{
	title:  "Best Selling Title",
	author: "Best Selling Author",
}

func main() {
	fmt.Println(GetBookStoreName())
	// "Go's Book Store"

	book := &Book{}
	book.SetTitle("A New Title")
	book.SetAuthor("A New Author")
	fmt.Println(book.GetTitle(), book.GetAuthor())
	// "A New Title A New Author"

	for k, v := range GetBooks() {
		fmt.Println(k, v.GetTitle(), v.GetAuthor())
		// "book 1 title 1 author 1"
		// "book 2 title 2 author 2"
		// "book 3 title 3 author 3"
	}

	fmt.Println(GetBestSellingAuthor())
	// "Best Selling Author"
}
