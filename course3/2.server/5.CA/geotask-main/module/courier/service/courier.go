package service

import (
	"context"
	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/module/courier/models"
	"gitlab.com/ptflp/geotask/module/courier/storage"
	"log"
	"math"
	"time"
)

// Направления движения курьера
const (
	DirectionUp    = 0
	DirectionDown  = 1
	DirectionLeft  = 2
	DirectionRight = 3
)

const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)

type Courierer interface {
	GetCourier(ctx context.Context) (*models.Courier, error)
	MoveCourier(courier models.Courier, direction, zoom int) error
}

type CourierService struct {
	courierStorage storage.CourierStorager
	allowedZone    geo.PolygonChecker
	disabledZones  []geo.PolygonChecker
}

func NewCourierService(courierStorage storage.CourierStorager, allowedZone geo.PolygonChecker, disbledZones []geo.PolygonChecker) Courierer {
	return &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZones: disbledZones}
}

func (c *CourierService) GetCourier(ctx context.Context) (*models.Courier, error) {
	// получаем курьера из хранилища используя метод GetOne из storage/courier.go
	one, err := c.courierStorage.GetOne(ctx)
	if err != nil {
		return &models.Courier{}, err
	}
	// проверяем, что курьер находится в разрешенной зоне
	if !c.allowedZone.Contains(geo.Point(one.Location)) {
		// если нет, то перемещаем его в случайную точку в разрешенной зоне
		NewLocation := c.allowedZone.RandomPoint()
		one.Location = models.Point(NewLocation)
		// сохраняем новые координаты курьера
		err = c.courierStorage.Save(ctx, *one)
		if err != nil {
			return &models.Courier{}, err
		}
	}
	return one, nil
}

//// MoveCourier : direction - направление движения курьера, zoom - зум карты
//func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
//	// точность перемещения зависит от зума карты использовать формулу 0.001 / 2^(zoom - 14)
//	accuracy := 0.001 / math.Pow(2, float64(zoom-14))
//	// 14 - это максимальный зум карты
//	newLng := courier.Location.Lng + float64(direction)*accuracy
//	newLat := courier.Location.Lat + float64(direction)*accuracy
//
//	// Создаем новую точку с новыми координатами
//	newLocation := models.Point{Lng: newLng, Lat: newLat}
//	// далее нужно проверить, что курьер не вышел за границы зоны
//	// если вышел, то нужно переместить его в случайную точку внутри зоны
//	if !c.allowedZone.Contains(geo.Point(newLocation)) {
//		// Если курьер вышел за границы зоны, перемещаем его в случайную точку внутри зоны
//		newLocation = models.Point(c.allowedZone.RandomPoint())
//	}
//	// сохраняем новые координаты курьера
//	ctx, _ := context.WithTimeout(context.Background(), time.Second)
//	err := c.courierStorage.Save(ctx, courier)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// MoveCourier : direction - направление движения курьера, zoom - зум карты
func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
	// точность перемещения зависит от зума карты использовать формулу 0.001 / 2^(zoom - 14)
	accuracy := 0.001 / math.Pow(2, float64(zoom-14))
	log.Println("accurancy: ", accuracy)
	log.Println("CourierService: сейчас координаты: ", courier.Location)
	// Изменяем координаты курьера в зависимости от направления
	switch direction {
	case DirectionUp:
		courier.Location.Lat += accuracy
	case DirectionDown:
		courier.Location.Lat -= accuracy
	case DirectionLeft:
		courier.Location.Lng -= accuracy
	case DirectionRight:
		courier.Location.Lng += accuracy
	}
	log.Println("CourierService: изменены координаты: ", courier.Location, " (", direction, ")")
	// Проверяем, что новые координаты находятся внутри разрешенной зоны
	if !c.allowedZone.Contains(geo.Point(courier.Location)) {
		// Если курьер вышел за границы зоны, перемещаем его в случайную точку внутри зоны
		newLocation := c.allowedZone.RandomPoint()
		courier.Location = models.Point(newLocation)
		log.Println("CourierService: назначены случайные координаты: ", courier.Location)
	}

	// Сохраняем новые координаты курьера
	log.Println("CourierService: отправлены координаты в редис: ", courier.Location)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := c.courierStorage.Save(ctx, courier); err != nil {
		return err
	}
	return nil
}

//func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
//	// Рассчитываем точность перемещения в зависимости от зума карты
//	accuracy := 0.001 / math.Pow(2, float64(zoom-14))
//
//	// Рассчитываем новые координаты курьера в зависимости от направления движения и точности перемещения
//	newLng := courier.Location.Lng + float64(direction)*accuracy
//	newLat := courier.Location.Lat + float64(direction)*accuracy
//
//	// Создаем новую точку с новыми координатами
//	newLocation := models.Point{Lng: newLng, Lat: newLat}
//
//	// Проверяем, что курьер не вышел за границы зоны
//	if !c.allowedZone.Contains(newLocation) {
//		// Если курьер вышел за границы зоны, перемещаем его в случайную точку внутри зоны
//		newLocation = c.allowedZone.RandomPoint()
//	}
//
//	// Обновляем координаты курьера
//	courier.Location = newLocation
//
//	// Сохраняем изменения в хранилище
//	err := c.storage.Update(ctx, courier)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
