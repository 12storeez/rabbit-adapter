package main

import (
	"context"
	"github.com/12storeez/pkg-go/rabbit"
	"github.com/gofiber/fiber"
	"github.com/zhs/loggr"
	"rabbit-adapter/config"
	"rabbit-adapter/services"
	"time"
)

func main() {
	// wait for rabbitmq start
	time.Sleep(60 * time.Second)

	cfg := config.New()
	logger := loggr.New("@version", cfg.App.Version)
	ctx := loggr.ToContext(context.TODO(), logger)

	logger.Info("Starting...")

	// connect to rabbitmq with reconnect
	rabbitConn, err := rabbit.NewConnection(ctx, cfg.Rabbit.Url)
	if err != nil {
		logger.Error(err)
	}

	// init stock publisher
	publisher := services.NewPublishStock(ctx, rabbitConn)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) {
		c.Next()
	})

	app.Post("/stocks", publisher.Publish())

	if err := app.Listen(cfg.App.Port); err != nil {
		logger.Fatal("can't run http server")
	}
}
