package storage

import (
	"context"
	"encoding/json"
	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"gitlab.com/ptflp/geotask/module/courier/models"
)

type CourierStorager interface {
	Save(ctx context.Context, courier models.Courier) error // сохранить курьера по ключу courier
	GetOne(ctx context.Context) (*models.Courier, error)    // получить курьера по ключу courier
}

type CourierStorage struct {
	storage *redis.Client
}

func (c *CourierStorage) Save(ctx context.Context, courier models.Courier) error {
	// Сериализуем курьера в формат JSON
	courierJSON, err := json.Marshal(courier)
	if err != nil {
		return err
	}

	// Сохраняем данные курьера в Redis по ключу "courier"
	err = c.storage.Set(ctx, "courier", string(courierJSON), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *CourierStorage) GetOne(ctx context.Context) (*models.Courier, error) {
	// Получаем данные о курьере из Redis по ключу "courier"
	courierJSON, err := c.storage.Get(ctx, "courier").Bytes()
	if err == redis.Nil {
		// Если курьер не найден, возвращаем nil, nil (нет данных о курьере)
		courierDef := models.Courier{10, models.Point{Lat: DefaultCourierLat, Lng: DefaultCourierLng}}
		return &courierDef, nil
	} else if err != nil {
		// Если произошла другая ошибка, возвращаем ее
		return nil, err
	}

	// Декодируем данные о курьере из JSON
	var courier models.Courier
	err = json.Unmarshal(courierJSON, &courier)
	if err != nil {
		// Если произошла ошибка при декодировании JSON, возвращаем ее
		return nil, err
	}

	// Возвращаем курьера и nil (без ошибки)
	return &courier, nil
}

func NewCourierStorage(storage *redis.Client) CourierStorager {
	return &CourierStorage{storage: storage}
}

// скопировал из courier
const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)
