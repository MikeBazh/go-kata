package model

type Author struct {
	ID    int
	Name  string
	Books []Book
}

type Book struct {
	ID    int
	Title string
}
