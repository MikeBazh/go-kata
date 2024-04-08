package orderModel

import "time"

type Order struct {
	Complete bool      `json:"complete"`
	Id       int       `json:"id"`
	PetId    int       `json:"petId"`
	Quantity int       `json:"quantity"`
	ShipDate time.Time `json:"shipDate"`
	Status   string    `json:"status"`
}

type Props map[string]int
