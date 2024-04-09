package controller

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.com/ptflp/geotask/module/courierfacade/service"
	"log"
	"net/http"
	"time"
)

type CourierController struct {
	courierService service.CourierFacer
}

func NewCourierController(courierService service.CourierFacer) *CourierController {
	return &CourierController{courierService: courierService}
}

type Courier struct {
	Score    int   `json:"score"`
	Location Point `json:"location"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (c *CourierController) GetStatus(ctx *gin.Context) {
	//log.Println("here")
	// Установить задержку в 50 миллисекунд
	time.Sleep(500 * time.Millisecond)

	// Получить статус курьера из сервиса courierService, используя метод GetStatus
	status := c.courierService.GetStatus(ctx)
	//resp := Courier{Location: Point{Lat: status.Courier.Location.Lat, Lng: status.Courier.Location.Lng}}
	log.Println("CourierController, получение/отправка статуса:", status.Courier.Location)

	// Подготовить данные для отправки в формате JSON
	responseData := map[string]interface{}{
		"courier": map[string]interface{}{
			"location": map[string]float64{
				"lat": status.Courier.Location.Lat,
				"lng": status.Courier.Location.Lng,
			},
		},
	}

	// Отправить статус курьера в ответ
	//ctx.JSON(http.StatusOK, gin.H{"status": status})
	ctx.JSON(http.StatusOK, responseData)
}

func (c *CourierController) MoveCourier(m webSocketMessage) {
	log.Println("im here")
	var cm CourierMove
	var err error

	// Получить данные из m.Data и десериализовать их в структуру CourierMove
	data, ok := m.Data.([]byte)

	log.Println("Получено сообщение:", data)

	if !ok {
		// Обработать ошибку, если есть
		log.Println("Invalid request data")
		return
	}
	err = json.Unmarshal(data, &cm)
	log.Println("Получено сообщение:", cm)

	if err != nil {
		// Обработать ошибку, если есть
		log.Println("Invalid request data")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	// Вызвать метод MoveCourier у courierService с необходимыми аргументами
	c.courierService.MoveCourier(ctx, cm.Direction, cm.Zoom)
	if err != nil {
		// Если произошла ошибка при перемещении курьера, вернуть ошибку
		log.Println("Failed to move courier:", err)
		return
	}

	// Отправить ответ об успешном перемещении курьера

	log.Println("Courier moved successfully")
}
