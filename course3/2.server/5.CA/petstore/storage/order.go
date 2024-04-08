package storage

import (
	"fmt"
	orderModel "go-kata/2.server/5.CA/petstore/dto/order"
)

func (s *LibraryStorage) Inventory() (props orderModel.Props, err error) {
	db, err := CreateTables()
	props = make(map[string]int)
	if err != nil {
		return orderModel.Props{}, err
	}
	// Выполняем запрос к базе данных для обновления
	row, err := db.Query("SELECT id, complete, petID, status, quantity, shipDate FROM orders")
	if err != nil {
		fmt.Println(err)
	}
	var order orderModel.Order
	for row.Next() {
		err = row.Scan(&order.Id, &order.Complete, &order.PetId, &order.Status, &order.Quantity, &order.ShipDate)
		if err != nil {
			return orderModel.Props{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		props[order.Status] = order.Quantity
	}
	return props, nil
}

func (s *LibraryStorage) AddOrder(order orderModel.Order) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO orders (petID, quantity, shipDate, status, complete) VALUES ($1, $2, $3, $4, $5)",
		order.PetId, order.Quantity, order.ShipDate, order.Status, order.Complete)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении order: %v", err)
	}
	return nil
}

func (s *LibraryStorage) FindOrderById(id int) (order orderModel.Order, err error) {
	db, err := CreateTables()
	if err != nil {
		return orderModel.Order{}, err
	}
	// Выполняем запрос к базе данных для обновления
	row, err := db.Query("SELECT id, complete, petID, status, quantity, shipDate FROM orders")
	for row.Next() {
		err = row.Scan(&order.Id, &order.Complete, &order.PetId, &order.Status, &order.Quantity, &order.ShipDate)
		if err != nil {
			return orderModel.Order{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		//fmt.Println(order)
	}
	return order, nil
}

func (s *LibraryStorage) DeleteOrder(id int) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для обновления
	var ID int
	query := "DELETE FROM orders WHERE id = $1 RETURNING id"
	row := db.QueryRow(query, id)
	err = row.Scan(&ID)
	if err != nil {
		return err
	}
	if ID == 0 {
		return fmt.Errorf("not found")
	}
	fmt.Println("order deleted")
	return nil
}
