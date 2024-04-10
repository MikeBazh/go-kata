package order

import (
	"context"
	"gitlab.com/ptflp/geotask/module/order/service"
	"log"
	"time"
)

const (
	// order generation interval
	orderGenerationInterval = 10 * time.Millisecond
	maxOrdersCount          = 200
)

// worker generates orders and put them into redis
type OrderGenerator struct {
	orderService service.Orderer
}

func NewOrderGenerator(orderService service.Orderer) *OrderGenerator {
	return &OrderGenerator{orderService: orderService}
}

func (o *OrderGenerator) Run() {
	go func() {
		ticker := time.NewTicker(orderGenerationInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				count, err := o.orderService.GetCount(ctx)
				//log.Println("OrderGenerator: count, err: ", count, err)
				if err != nil {
					log.Printf("Ошибка при получении количества заказов: %v", err)
					continue
				}
				//log.Println("Заказ сгенерирован", err)

				if count >= maxOrdersCount {
					continue
				}

				if err := o.orderService.GenerateOrder(ctx); err != nil {
					log.Printf("Ошибка при генерации заказа: %v", err)
					continue
				}
				//log.Println("Заказ сгенерирован", err)
			}
		}
	}()
}

// запускаем горутину, которая будет генерировать заказы не более чем раз в 10 миллисекунд
// не более 200 заказов используя константы orderGenerationInterval и maxOrdersCount
// нужно использовать метод orderService.GetCount() для получения количества заказов
// и метод orderService.GenerateOrder() для генерации заказа
// если количество заказов меньше maxOrdersCount, то нужно сгенерировать новый заказ
// если количество заказов больше или равно maxOrdersCount, то не нужно ничего делать
// если при генерации заказа произошла ошибка, то нужно вывести ее в лог
// если при получении количества заказов произошла ошибка, то нужно вывести ее в лог
// внутри горутины нужно использовать select и time.NewTicker()

//for i := 1; i <= maxOrdersCount; i++ {
//err := o.orderService.GenerateOrder(ctx)
//if err != nil {
//return
//}
//}
//time.Sleep(orderGenerationInterval)
//go func() {
//	time.Sleep(orderGenerationInterval)
//	num, err := o.orderService.GetCount(ctx)
//	if err != nil {
//		log.Println(err)
//	}
//	if num < maxOrdersCount {
//		err := o.orderService.GenerateOrder(ctx)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}()
