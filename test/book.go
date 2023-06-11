package test

//go:generate goaccessor -s bookStoreName -g
const bookStoreName = "Go's Book Store"

// TODO go:generate goaccessor --symbol Book --getter --setter
type Book struct {
	Title  string
	Author string
}

// TODO go:generate goaccessor --symbol books --getter --setter
// var books = map[string]*Book {
//     ...
// }

func main() {}
