package services

import (
	"github.com/gofiber/fiber"
)

type Publisher interface {
	Publish() func(*fiber.Ctx)
}
