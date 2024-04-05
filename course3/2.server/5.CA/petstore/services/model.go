package services

import "time"

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	UserStatus int    `json:"userStatus"`
}

type Order struct {
	Complete bool      `json:"complete"`
	Id       int       `json:"id"`
	PetId    int       `json:"petId"`
	Quantity int       `json:"quantity"`
	ShipDate time.Time `json:"shipDate"`
	Status   string    `json:"status"`
}

type Props map[string]int

type Pet struct {
	Category struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status    string   `json:"status"`
	Tags      []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}
