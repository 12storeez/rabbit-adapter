package services

import (
	"context"
	"encoding/json"
	"github.com/12storeez/pkg-go/rabbit"
	"github.com/gofiber/fiber"
	"github.com/zhs/loggr"
	"net/http"
	"rabbit-adapter/models"
)

type PublishStockPayload struct {
	Exchange string
	Key string
	Data []StocksPayload
}

type StocksPayload struct {
	Barcode string
	StoreID int
	Available int
	Reserved int
}

func NewPublishStock(ctx context.Context, conn *rabbit.Connection) Publisher {
	return &publishStock{
		ctx: ctx,
		rabbitConn: conn,
		inch: make(chan PublishStockPayload),
		outch: make(chan PublishStockPayload),
	}
}

type publishStock struct {
	ctx context.Context
	rabbitConn *rabbit.Connection
	inch chan PublishStockPayload
	outch chan PublishStockPayload
	exchanges []string
}

func (p publishStock) Publish() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		log := loggr.WithContext(p.ctx).
			With("method", c.Method(),
				"service", "StockPublisher",
				"request", c.Body(),
			)

		log.Info("[request]")

		var body PublishStockPayload
		if err := c.BodyParser(&body); err != nil {
			log.Error("error parsing body: " + err.Error())
			c.Status(http.StatusBadRequest)
			_ = c.JSON(models.ErrResponse{
				Err: err.Error(),
			})
			return
		}

		p.declareExchange(body.Exchange)

		go func() {
			var in <-chan PublishStockPayload = p.inch
			var out chan<- PublishStockPayload
			var val PublishStockPayload
			for {
				select {
				case out <- val:
					out = nil
					in = p.inch
				case val = <-in:
					out = p.outch
					in = nil
				}
			}
		}()

		go func() {
			for payload := range p.outch {
				data, _ := json.Marshal(payload.Data)
				publish, _ := p.rabbitConn.NewPublisher(payload.Exchange, payload.Key)
				err := publish(data)
				if err != nil {
					log.Errorf("error publish to exchange: %s, key: %s. Error: %v", payload.Exchange, payload.Key, err)
				}
			}
		}()

		p.inch <- body

		resp := models.OkResponse{
			Ok: true,
		}
		_ = c.JSON(resp)
	}
}

func (p *publishStock) declareExchange(exchange string)  {
	for _, declared := range p.exchanges {
		if exchange == declared {
			return
		}
	}

	_ = p.rabbitConn.Channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)

	p.exchanges = append(p.exchanges, exchange)
}

