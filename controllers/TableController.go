package controllers

import (
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TableCollection = database.OpenCollection(database.Client, "table")

func GetTables(c *fiber.Ctx) error {
	tables := models.Table{}

	cursor, err := TableCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	if err := cursor.All(c.Context(), &tables); err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": tables})
}

func GetTable(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}

	tableId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "invalid tableId"})
	}

	table := models.Table{}
	err = TableCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: tableId}}).Decode(&table)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "cannot get the table"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "fail", "data": table})
}

func CreateTable(c *fiber.Ctx) error {
	table := models.Table{}

	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	validationErr := Validate.Struct(table)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": validationErr.Error()})
	}

	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertedItem, err := TableCollection.InsertOne(c.Context(), table)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "table creation failed"})
	}

	newTable := models.Table{}
	err = TableCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newTable)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "table creation failed"})
	}

	return c.Status(201).JSON(&fiber.Map{"status": "success", "data": newTable})
}

func UpdateTable(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}

	tableId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "invalid table id"})
	}

	table := models.Table{}
	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	valiationErr := Validate.Struct(table)
	if valiationErr != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": valiationErr.Error()})
	}

	tableObj := primitive.D{}

	if table.NumberOfGuests != nil {
		tableObj = append(tableObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
	}
	if table.TableNumber != nil {
		tableObj = append(tableObj, bson.E{Key: "table_number", Value: table.TableNumber})
	}

	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	tableObj = append(tableObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

	filter := bson.D{{Key: "_id", Value: tableId}}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := TableCollection.UpdateOne(c.Context(), filter, bson.D{{Key: "$set", Value: tableObj}}, &opt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "table update failed"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": result})
}
