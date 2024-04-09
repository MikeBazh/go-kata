package storage

import (
	"context"
	"encoding/json"
	"fmt"

	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"gitlab.com/ptflp/geotask/module/order/models"
	"strconv"
	"time"
)

type OrderStorager interface {
	Save(ctx context.Context, order models.Order, maxAge time.Duration) error                       // сохранить заказ с временем жизни
	GetByID(ctx context.Context, orderID int) (*models.Order, error)                                // получить заказ по id
	GenerateUniqueID(ctx context.Context) (int64, error)                                            // сгенерировать уникальный id
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // получить заказы в радиусе от точки
	GetCount(ctx context.Context) (int, error)                                                      // получить количество заказов
	RemoveOldOrders(ctx context.Context, maxAge time.Duration) error                                // удалить старые заказы по истечению времени maxAge
}

type OrderStorage struct {
	storage *redis.Client
}

func NewOrderStorage(storage *redis.Client) OrderStorager {
	return &OrderStorage{storage: storage}
}

func (o *OrderStorage) Save(ctx context.Context, order models.Order, maxAge time.Duration) error {
	// save with geo redis
	return o.saveOrderWithGeo(ctx, order, maxAge)
}

func (o *OrderStorage) RemoveOldOrders(ctx context.Context, maxAge time.Duration) error {
	// получить ID всех старых ордеров, которые нужно удалить
	// используя метод ZRangeByScore
	// старые ордеры это те, которые были созданы две минуты назад
	// и более
	/**
	&redis.ZRangeBy{
		Max: использовать вычисление времени с помощью maxAge,
		Min: "0",
	}
	*/

	// Проверить количество старых ордеров
	// удалить старые ордеры из redis используя метод ZRemRangeByScore где ключ "orders" min "-inf" max "(время создания старого ордера)"
	// удалять ордера по ключу не нужно, они будут удалены автоматически по истечению времени жизни

	return nil
}

func (o *OrderStorage) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	// Получаем данные о заказе из Redis по ключу "order:ID"
	orderJSON, err := o.storage.Get(ctx, "order:"+strconv.Itoa(orderID)).Bytes()
	if err == redis.Nil {
		// Если заказ не найден, возвращаем nil, nil (нет данных о заказе)
		return nil, nil
	} else if err != nil {
		// Если произошла другая ошибка, возвращаем ее
		return nil, err
	}

	// Декодируем данные о заказе из JSON
	var order models.Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		// Если произошла ошибка при декодировании JSON, возвращаем ее
		return nil, err
	}

	// Возвращаем заказ и nil (без ошибки)
	return &order, nil
}

func (o *OrderStorage) saveOrderWithGeo(ctx context.Context, order models.Order, maxAge time.Duration) error {
	var err error

	// Serialize the order into JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Save the order in JSON format to Redis with a specific key
	key := "orderID:" + strconv.FormatInt(order.ID, 10)
	err = o.storage.Set(ctx, key, string(orderJSON), maxAge).Err()
	if err != nil {
		return err
	}

	// Add the order to the geospatial index using GeoAdd
	err = o.storage.GeoAdd(ctx, "orders_geo_index", &redis.GeoLocation{
		Name:      key,
		Longitude: order.Lng,
		Latitude:  order.Lat,
	}).Err()
	if err != nil {
		return err
	}

	// Save the order into a zset for easy retrieval based on the order creation time
	_, err = o.storage.ZAdd(ctx, "orders", &redis.Z{
		Score:  float64(order.CreatedAt.Unix()), // Use Unix timestamp as the score
		Member: key,
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderStorage) GetCount(ctx context.Context) (int, error) {
	// получить количество ордеров в упорядоченном множестве используя метод ZCard
	count, err := o.storage.ZCard(ctx, "orders").Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (o *OrderStorage) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	var err error
	var orders []models.Order
	var data []byte
	var ordersLocation []redis.GeoLocation

	// используем метод getOrdersByRadius для получения ID заказов в радиусе
	ordersLocation, err = o.getOrdersByRadius(ctx, lng, lat, radius, unit)
	// обратите внимание, что в случае отсутствия заказов в радиусе
	// метод getOrdersByRadius должен вернуть nil, nil (при ошибке redis.Nil)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	orders = make([]models.Order, 0, len(ordersLocation))
	// проходим по списку ID заказов и получаем данные о заказе
	for _, orderLocation := range ordersLocation {
		// получаем данные о заказе по ID из redis по ключу order:ID
		data, err = o.storage.Get(ctx, "orderID:"+orderLocation.Name).Bytes()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var order models.Order
		err = json.Unmarshal(data, &order)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Добавляем заказ в список
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *OrderStorage) getOrdersByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]redis.GeoLocation, error) {
	// в данном методе мы получаем список ордеров в радиусе от точки
	// возвращаем список ордеров с координатами и расстоянием до точки

	query := &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	}
	location, err := o.storage.GeoRadius(ctx, "orders_geo_index", lng, lat, query).Result()
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (o *OrderStorage) GenerateUniqueID(ctx context.Context) (int64, error) {
	// Используем INCR для генерации уникального идентификатора
	id, err := o.storage.Incr(ctx, "order:id").Result()
	if err != nil {
		return 0, err
	}

	return id, nil
}