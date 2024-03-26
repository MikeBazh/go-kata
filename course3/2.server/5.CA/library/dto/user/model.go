package UserModel

type User struct {
	ID          int
	Name        string
	RentedBooks []RentedBook
}

type RentedBook struct {
	ID       int
	Title    string
	Author   string
	IsRented bool
}

//type Author struct {
//	ID   int
//	Name string
//}
