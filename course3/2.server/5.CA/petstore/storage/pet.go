package storage

import (
	"encoding/json"
	"fmt"
	PetModel "go-kata/2.server/5.CA/petstore/dto/pet"
)

func (s *LibraryStorage) AddPet(pet PetModel.Pet) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для добавления
	//var tags []byte
	tags, err := json.Marshal(pet.Tags)
	category, err := json.Marshal(pet.Category)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO pets (name, status, tags, category) VALUES ($1, $2, $3, $4)", pet.Name, pet.Status, tags, category)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении pet: %v", err)
	}
	fmt.Println("pet added")
	return nil
}

func (s *LibraryStorage) UpdatePet(pet PetModel.Pet) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	tags, err := json.Marshal(pet.Tags)
	category, err := json.Marshal(pet.Category)
	// Выполняем запрос к базе данных для обновления
	_, err = db.Exec("UPDATE pets SET name=$1, status=$2, tags=$3, category=$4 WHERE id=$5",
		pet.Name, pet.Status, tags, category, pet.Id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении pet: %v", err)
	}
	fmt.Println("pet updated")
	return nil
}

func (s *LibraryStorage) UpdatePetWithData(pet PetModel.Pet) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для обновления
	_, err = db.Exec("UPDATE pets SET name=$1, status=$2 WHERE id=$3",
		pet.Name, pet.Status, pet.Id)
	//fmt.Println(id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении pet: %v", err)
	}
	return nil
}

func (s *LibraryStorage) FindPetById(id int) (pet PetModel.Pet, err error) {
	db, err := CreateTables()
	if err != nil {
		return PetModel.Pet{}, err
	}
	// Выполняем запрос к базе данных для обновления
	row, err := db.Query("SELECT * FROM pets WHERE id=$1", id)
	var tags Tags
	var category Category
	var tagsJson []byte
	var categoryJson []byte
	for row.Next() {
		err = row.Scan(&pet.Id, &pet.Name, &pet.Status, &tagsJson, &categoryJson)
		err = json.Unmarshal(tagsJson, &tags)
		err = json.Unmarshal(categoryJson, &category)
		if err != nil {
			return PetModel.Pet{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		pet.Tags = tags
		pet.Category = category
		//fmt.Println(pet)
	}
	return pet, nil
}

func (s *LibraryStorage) DeletePet(id int) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для обновления
	var ID int
	query := "DELETE FROM pets WHERE id = $1 RETURNING id"
	row := db.QueryRow(query, id)
	err = row.Scan(&ID)
	if err != nil {
		return err
	}
	if ID == 0 {
		return fmt.Errorf("not found")
	}
	fmt.Println("pet deleted")
	return nil
}

func (s *LibraryStorage) FindPetByStatus(status string) (pets []PetModel.Pet, err error) {
	db, err := CreateTables()
	if err != nil {
		return []PetModel.Pet{}, err
	}
	var pet PetModel.Pet
	// Выполняем запрос к базе данных для обновления
	rows, err := db.Query("SELECT * FROM pets WHERE status=$1", status)
	if err != nil {
		return []PetModel.Pet{}, fmt.Errorf("ошибка при запросе pet: %v", err)
	}
	for rows.Next() {
		var tags Tags
		var category Category
		var tagsJson []byte
		var categoryJson []byte
		err := rows.Scan(&pet.Id, &pet.Name, &pet.Status, &tagsJson, &categoryJson)
		err = json.Unmarshal(tagsJson, &tags)
		err = json.Unmarshal(categoryJson, &category)
		if err != nil {
			return []PetModel.Pet{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		pet.Tags = tags
		pet.Category = category
		//fmt.Println(pet)
		pets = append(pets, pet)
	}
	return pets, nil
}

type Tags []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
