package storage

import (
	"fmt"
	_ "github.com/lib/pq"
	orderModel "go-kata/2.server/5.CA/petstore/dto/order"
	PetModel "go-kata/2.server/5.CA/petstore/dto/pet"

	//UserModel "go-kata/2.server/5.CA/library/dto/user"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
)

type LibraryRepository interface {
	CreateUser(UserModel.User) error
	CreateWithArray([]UserModel.User) error
	CreateWithList([]UserModel.User) error
	GetUserByName(name string) (UserModel.User, error)
	UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error)
	DeleteUserByName(name string) (UserModel.User, error)
	LoginUser(name, password string) error
	LogoutUser(name string) error
	//
	FindPetByStatus(status string) (pets []PetModel.Pet, err error)
	AddPet(pet PetModel.Pet) error
	UpdatePet(pet PetModel.Pet) error
	UpdatePetWithData(pet PetModel.Pet) error
	FindPetById(int) (pet PetModel.Pet, err error)
	DeletePet(id int) error
	//
	Inventory() (props orderModel.Props, err error)
	AddOrder(order orderModel.Order) error
	FindOrderById(id int) (order orderModel.Order, err error)
	DeleteOrder(id int) error
}

type LibraryStorage struct {
}

const (
	connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
	//connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
)

// NewLibraryStorage - конструктор хранилища пользователей
func NewLibraryStorage() *LibraryStorage {
	return &LibraryStorage{}
}

func (ls *LibraryStorage) GetUserByName(name string) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	var user UserModel.User
	query := "SELECT id, username, firstname, lastname, email, phone, userStatus FROM users WHERE username = $1"
	row := db.QueryRow(query, name)
	err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
	user.Password = "*****"
	if err != nil {
		return UserModel.User{}, err
	}
	return user, nil
}

func (ls *LibraryStorage) UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	//var updatedUser UserModel.User
	query := "UPDATE users SET firstname = $1, lastname=$2, email = $3, phone = $4, userStatus=$5, password=$6  WHERE username = $7 RETURNING id"
	var id int
	err = db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.UserStatus, newUser.Password, name).Scan(&id)
	if err != nil {
		return UserModel.User{}, err
	}
	fmt.Println("user", id, "updated")
	return newUser, nil
}

func (ls *LibraryStorage) DeleteUserByName(name string) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	var deletedUser UserModel.User
	query := "DELETE FROM users WHERE username = $1 RETURNING id, username"
	row := db.QueryRow(query, name)
	err = row.Scan(&deletedUser.Id, &deletedUser.Username)
	if err != nil {
		return UserModel.User{}, err
	}
	fmt.Println("user", deletedUser.Username, "deleted")
	return deletedUser, nil
}

func (ls *LibraryStorage) LoginUser(name, password string) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	//m:=make(map[string]string)
	var UserPassword string
	query := "SELECT password FROM users WHERE username = $1"
	err = db.QueryRow(query, name).Scan(&UserPassword)
	if err != nil {
		return err
	}
	var ID string
	if UserPassword == password {
		// Implement login logic here
		query := "UPDATE users SET userStatus = 1 WHERE username = $1 RETURNING id"
		row := db.QueryRow(query, name)
		err = row.Scan(&ID)
		if err != nil {
			return err
		}
		fmt.Println("user ", ID, "logged in")
	} else {
		fmt.Println("user не найден или пароль не совпадает")
		{
			return fmt.Errorf("user не найден или пароль не совпадает")
		}
	}
	return nil
}

func (ls *LibraryStorage) LogoutUser(name string) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	var ID string
	query := "UPDATE users SET userStatus = 0 WHERE username = $1 RETURNING id"
	row := db.QueryRow(query, name)
	err = row.Scan(&ID)
	if err != nil {
		return err
	}
	fmt.Println("user ", ID, "logged out")
	return nil
}

func (ls *LibraryStorage) CreateWithArray(users []UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	for _, user := range users {
		query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
		//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ls *LibraryStorage) CreateWithList(users []UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	for _, user := range users {
		query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
		//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ls *LibraryStorage) CreateUser(user UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
	//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
	if err != nil {
		return err
	}
	return nil
}
