package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func GetOrderItem(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func GetOrderItemsByOrderId(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func ItemsByOrderId(id string) (OrderItems []primitive.M, err error) {
	return
}

func CreateOrderItem(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func UpdateOrderItem(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
