package test

//go:generate goaccessor -t bookStoreName -g
const bookStoreName = "Go's Book Store"

// TODO go:generate goaccessor --target Book --getter --setter
type Book struct {
	Title  string
	Author string
}

// TODO go:generate goaccessor --target books --getter --setter
// var books = map[string]*Book {
//     ...
// }

func main() {}
