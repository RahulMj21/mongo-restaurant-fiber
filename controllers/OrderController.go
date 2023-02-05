package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var OrderCollection = database.OpenCollection(database.Client, "order")

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}

	cursor, err := OrderCollection.Find(c.Context(), bson.D{{}})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	err = cursor.All(c.Context(), &orders)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": orders})
}

func GetOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}

	orderId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	order := models.Order{}
	if err := OrderCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: orderId}}).Decode(&order); err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "order not found"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": order})
}

func CreateOrder(c *fiber.Ctx) error {
	order := models.Order{}
	table := models.Table{}
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}
	validationErr := Validate.Struct(order)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": validationErr.Error()})
	}

	if order.TableId != nil {
		tableId, err := primitive.ObjectIDFromHex(*order.TableId)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "invalid table_id"})
		}
		if err := TableCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: tableId}}).Decode(&table); err != nil {
			msg := fmt.Sprintf("table not found with the table_id: %s", *order.TableId)
			return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": msg})
		}
	}

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertedItem, err := OrderCollection.InsertOne(c.Context(), order)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "failed to create new order"})
	}

	newOrder := models.Order{}
	err = OrderCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newOrder)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "failed to create new order"})
	}

	return c.Status(201).JSON(&fiber.Map{"status": "success", "data": newOrder})
}

func UpdateOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}

	orderId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	filter := bson.D{{Key: "_id", Value: orderId}}

	order := models.Order{}
	table := models.Table{}

	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	orderObj := primitive.D{}

	if order.TableId != nil {
		tableId, err := primitive.ObjectIDFromHex(*order.TableId)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
		}

		if err := TableCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: tableId}}).Decode(&table); err != nil {
			return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
		}
		orderObj = append(orderObj, bson.E{Key: "table_id", Value: order.TableId})
	}

	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderObj = append(orderObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := OrderCollection.UpdateOne(c.Context(), filter, bson.D{{Key: "$set", Value: orderObj}}, &opt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": result})
}

func OrderItemOrderCreater(order models.Order) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	insertedItem, err := OrderCollection.InsertOne(ctx, order)

	if err != nil {
		return insertedItem.InsertedID, err
	}
	return insertedItem.InsertedID, nil
}
