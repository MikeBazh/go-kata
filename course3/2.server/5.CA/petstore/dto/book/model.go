package BookModel

type Book struct {
	ID     int
	Title  string
	Author Author
}

type Author struct {
	ID   int
	Name string
}
