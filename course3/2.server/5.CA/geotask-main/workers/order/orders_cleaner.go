package order

import (
	"context"
	"gitlab.com/ptflp/geotask/module/order/service"
	"log"
	"time"
)

const (
	orderCleanInterval = 5 * time.Second
)

// OrderCleaner воркер, который удаляет старые заказы
// используя метод orderService.RemoveOldOrders()
type OrderCleaner struct {
	orderService service.Orderer
}

func NewOrderCleaner(orderService service.Orderer) *OrderCleaner {
	return &OrderCleaner{orderService: orderService}
}

func (o *OrderCleaner) Run() {
	//ticker := time.NewTicker(orderCleanInterval)
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer ticker.Stop()
	go func() {
		ticker := time.NewTicker(orderCleanInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				err := o.orderService.RemoveOldOrders(ctx)
				if err != nil {
					log.Printf("Ошибка при удалении старых заказов: %v", err)
					continue
				}
			}
		}
	}()
	// исользовать горутину и select
	// внутри горутины нужно использовать time.NewTicker()
	// и вызывать метод orderService.RemoveOldOrders()
	// если при удалении заказов произошла ошибка, то нужно вывести ее в лог
}
