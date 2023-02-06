package controllers

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var OrderItemsCollection = database.OpenCollection(database.Client, "order_items")

type OrderItemPack struct {
	TableId    string             `json:"table_id"`
	OrderItems []models.OrderItem `json:"order_items"`
}

func GetOrderItems(c *fiber.Ctx) error {
	filter := bson.D{{}}
	cursor, err := OrderItemsCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	orderItems := []models.OrderItem{}
	if err := cursor.All(c.Context(), &orderItems); err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": orderItems})
}

func GetOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}
	orderItemId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	orderItem := models.OrderItem{}
	err = OrderItemsCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: orderItemId}}).Decode(&orderItem)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{"status": "fail", "message": "order_item not found"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": orderItem})
}

func GetOrderItemsByOrderId(c *fiber.Ctx) error {
	orderIdParam := c.Params("order_id")
	if orderIdParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "order_id cannot be empty"})
	}
	allOrderItems, err := ItemsByOrderId(orderIdParam)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": allOrderItems})
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
