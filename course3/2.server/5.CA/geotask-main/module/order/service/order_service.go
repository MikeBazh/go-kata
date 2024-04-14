package service

import (
	"context"
	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/module/order/models"
	"gitlab.com/ptflp/geotask/module/order/storage"
	"math/rand"
	"time"
)

const (
	minDeliveryPrice = 100.00
	maxDeliveryPrice = 500.00

	maxOrderPrice = 3000.00
	minOrderPrice = 1000.00

	orderMaxAge = 2 * time.Minute
)

type Orderer interface {
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // возвращает заказы через метод storage.GetByRadius
	Save(ctx context.Context, order models.Order) error                                             // сохраняет заказ через метод storage.Save с заданным временем жизни OrderMaxAge
	GetCount(ctx context.Context) (int, error)                                                      // возвращает количество заказов через метод storage.GetCount
	RemoveOldOrders(ctx context.Context) error                                                      // удаляет старые заказы через метод storage.RemoveOldOrders с заданным временем жизни OrderMaxAge
	GenerateOrder(ctx context.Context) error                                                        // генерирует заказ в случайной точке из разрешенной зоны, с уникальным id, ценой и ценой доставки
}

// OrderService реализация интерфейса Orderer
// в нем должны быть методы GetByRadius, Save, GetCount, RemoveOldOrders, GenerateOrder
// данный сервис отвечает за работу с заказами
type OrderService struct {
	storage       storage.OrderStorager
	allowedZone   geo.PolygonChecker
	disabledZones []geo.PolygonChecker
}

func (o *OrderService) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	order, err := o.storage.GetByRadius(ctx, lng, lat, radius, unit)
	//log.Println("OrderService: ", order, err)
	if err != nil {
		return nil, err
	}
	return order, err
}

func (o *OrderService) Save(ctx context.Context, order models.Order) error {
	err := o.storage.Save(ctx, order, 0)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderService) GetCount(ctx context.Context) (int, error) {
	count, err := o.storage.GetCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (o *OrderService) RemoveOldOrders(ctx context.Context) error {
	err := o.storage.RemoveOldOrders(ctx, orderMaxAge)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderService) GenerateOrder(ctx context.Context) error {
	// Генерация заказа в случайной точке из разрешенной зоны
	//randomPoint := o.allowedZone.RandomPoint()
	randomPoint := geo.GetRandomAllowedLocation(o.allowedZone, o.disabledZones)
	//log.Println("OrderService: GenerateOrder: randomPoint", randomPoint)
	orderID, err := o.storage.GenerateUniqueID(ctx)
	//log.Println("OrderService: GenerateOrder: orderID, err", orderID, err)
	if err != nil {
		return err
	}
	rand.Seed(1)
	order := models.Order{
		ID:            orderID,
		Price:         rand.Float64()*(maxOrderPrice-minOrderPrice) + minOrderPrice,
		DeliveryPrice: rand.Float64()*(maxDeliveryPrice-minDeliveryPrice) + minDeliveryPrice,
		Lng:           randomPoint.Lng,
		Lat:           randomPoint.Lat,
		IsDelivered:   false,
		CreatedAt:     time.Now(),
	}
	//log.Println("OrderService: GenerateOrder: ", order)
	// Сохранение заказа
	err = o.storage.Save(ctx, order, orderMaxAge)
	if err != nil {
		return err
	}
	return nil
}

func NewOrderService(storage storage.OrderStorager, allowedZone geo.PolygonChecker, disallowedZone []geo.PolygonChecker) Orderer {
	return &OrderService{storage: storage, allowedZone: allowedZone, disabledZones: disallowedZone}
}
