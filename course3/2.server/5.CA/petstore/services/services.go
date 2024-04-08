package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	orderModel "go-kata/2.server/5.CA/petstore/dto/order"
	PetModel "go-kata/2.server/5.CA/petstore/dto/pet"
	//model "go-kata/2.server/5.CA/library/dto"
	//BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
	"go-kata/2.server/5.CA/petstore/storage"
)

var users = make(map[string]string)
var tokens = make(map[string]bool)

//var tb = TokenBlacklist{tokens: tokens}

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	CreateUser(UserModel.User) error
	GetUserByName(name string) (UserModel.User, error)
	UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error)
	DeleteUserByName(name string) (UserModel.User, error)
	LoginUser(name, password string) (string, error)
	LogoutUser(name string) error
	CreateWithArray([]UserModel.User) error
	CreateWithList([]UserModel.User) error
	//
	FindPetByStatus(status string) (pets []PetModel.Pet, err error)
	AddPet(pet PetModel.Pet) error
	UpdatePet(pet PetModel.Pet) error
	FindPetById(int) (pet PetModel.Pet, err error)
	UpdatePetWithData(id int, name, status string) error
	DeletePet(id int) error
	//
	Inventory() (props orderModel.Props, err error)
	FindOrderById(id int) (order orderModel.Order, err error)
	DeleteOrder(id int) error
	AddOrder(order orderModel.Order) error
}

type Service struct {
	UserStorage storage.LibraryRepository
}

func NewService(UserStorage storage.LibraryRepository) *Service {
	return &Service{
		UserStorage: UserStorage}
}

func (s *Service) Inventory() (props orderModel.Props, err error) {
	props, err = s.UserStorage.Inventory()
	if err != nil {
		fmt.Println(err)
		return props, err
	}
	return props, nil
}

func (s *Service) FindOrderById(id int) (order orderModel.Order, err error) {
	order, err = s.UserStorage.FindOrderById(id)
	if err != nil {
		fmt.Println(err)
		return order, err
	}
	return order, nil
}

func (s *Service) DeleteOrder(id int) error {
	err := s.UserStorage.DeleteOrder(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) AddOrder(order orderModel.Order) error {
	err := s.UserStorage.AddOrder(order)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) CreateUser(user UserModel.User) error {
	err := s.UserStorage.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetUserByName(name string) (UserModel.User, error) {
	User, err := s.UserStorage.GetUserByName(name)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error) {
	User, err := s.UserStorage.UpdateUserByName(name, newUser)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) DeleteUserByName(name string) (UserModel.User, error) {
	User, err := s.UserStorage.DeleteUserByName(name)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) LoginUser(name, password string) (string, error) {
	err := s.UserStorage.LoginUser(name, password)
	var tokenString string
	if err != nil {
		fmt.Println("service:", err)
		return "", err
	} else {
		_, tokenString, err = TokenAuth.Encode(jwt.MapClaims{"sub": name})
		if err != nil {
			return "", fmt.Errorf("ошибка генерации JWT токена")
		}

	}
	return tokenString, nil
}

//type TokenBlacklist struct {
//	sync.Mutex
//	tokens map[string]bool // map для быстрого поиска токенов
//}
//
//// Функция для добавления токена в черный список
//func (tb *TokenBlacklist) Add(token string) {
//	tb.Lock()
//	defer tb.Unlock()
//	tb.tokens[token] = true
//}
//
//// Функция для проверки, есть ли токен в черном списке
//func (tb *TokenBlacklist) Contains(token string) bool {
//	tb.Lock()
//	defer tb.Unlock()
//	return tb.tokens[token]
//}

func (s *Service) LogoutUser(name string) error {
	err := s.UserStorage.LogoutUser(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) FindPetByStatus(status string) (pets []PetModel.Pet, err error) {
	pets, err = s.UserStorage.FindPetByStatus(status)
	if err != nil {
		fmt.Println("service:", err)
		return []PetModel.Pet{}, err
	}
	return pets, nil
}

func (s *Service) AddPet(pet PetModel.Pet) error {
	err := s.UserStorage.AddPet(pet)
	if err != nil {
		fmt.Println("service:", err)
		return err
	}
	return nil
}

func (s *Service) UpdatePet(pet PetModel.Pet) error {
	err := s.UserStorage.UpdatePet(pet)
	if err != nil {
		fmt.Println("service:", err)
		return err
	}
	return nil
}

func (s *Service) FindPetById(id int) (pet PetModel.Pet, err error) {
	pet, err = s.UserStorage.FindPetById(id)
	if err != nil {
		fmt.Println("service:", err)
		return PetModel.Pet{}, err
	}
	return pet, nil
}

func (s *Service) UpdatePetWithData(id int, name, status string) error {
	pet := PetModel.Pet{Id: id, Name: name, Status: status}
	//fmt.Println(pet)
	err := s.UserStorage.UpdatePetWithData(pet)
	if err != nil {
		fmt.Println("service:", err)
		return err
	}
	return nil
}

func (s *Service) DeletePet(id int) error {
	err := s.UserStorage.DeletePet(id)
	if err != nil {
		fmt.Println("service:", err)
		return err
	}
	return nil
}

func (s *Service) CreateWithArray(users []UserModel.User) error {
	err := s.UserStorage.CreateWithArray(users)
	if err != nil {
		fmt.Println("Ошибка сервис:", err)
		return err
	}
	return nil
}

func (s *Service) CreateWithList(users []UserModel.User) error {
	err := s.UserStorage.CreateWithList(users)
	if err != nil {
		fmt.Println("Ошибка сервис:", err)
		return err
	}
	return nil
}
