package service

import (
	"context"
	cm "gitlab.com/ptflp/geotask/module/courier/models"
	cservice "gitlab.com/ptflp/geotask/module/courier/service"
	cfm "gitlab.com/ptflp/geotask/module/courierfacade/models"
	om "gitlab.com/ptflp/geotask/module/order/models"
	oservice "gitlab.com/ptflp/geotask/module/order/service"
	"log"
)

const (
	CourierVisibilityRadius = 2800 // 2.8km
)

type CourierFacer interface {
	MoveCourier(ctx context.Context, direction, zoom int) // отвечает за движение курьера по карте direction - направление движения, zoom - уровень зума
	GetStatus(ctx context.Context) cfm.CourierStatus      // отвечает за получение статуса курьера и заказов вокруг него
}

// CourierFacade фасад для курьера и заказов вокруг него (для фронта)
type CourierFacade struct {
	courierService cservice.Courierer
	orderService   oservice.Orderer
}

func (cf *CourierFacade) MoveCourier(ctx context.Context, direction, zoom int) {
	//log.SetOutput(os.Stdout)
	// Вызываем метод MoveCourier из courierService
	One, err := cf.courierService.GetCourier(ctx)
	log.Println("CourierFacade получил координаты: ", One.Location)
	_ = cf.courierService.MoveCourier(*One, direction, zoom)

	if err != nil {
		// Обработка ошибки
		log.Println(err)
	}
}

type CourierStatus struct {
	Courier cm.Courier `json:"courier"`
	Orders  []om.Order `json:"orders"`
}

func (cf *CourierFacade) GetStatus(ctx context.Context) cfm.CourierStatus {

	// Получаем статус курьера из courierService
	One, err := cf.courierService.GetCourier(ctx)
	log.Println("CourierFacade, получение/отправка статуса:", One.Location)
	// Получаем заказы вокруг курьера из orderService
	orders, err := cf.orderService.GetByRadius(ctx, One.Location.Lat, One.Location.Lng, CourierVisibilityRadius, "m")
	if err != nil {
		// Обработка ошибки
		log.Println(err)
	}

	// Создаем объект cfm.CourierStatus и заполняем его данными
	courierStatus := cfm.CourierStatus{
		Courier: cm.Courier{
			Score:    444,
			Location: One.Location,
		},
		Orders: orders,
	}
	log.Println(One.Location)
	return courierStatus
}

func NewCourierFacade(courierService cservice.Courierer, orderService oservice.Orderer) CourierFacer {
	return &CourierFacade{courierService: courierService, orderService: orderService}
}
