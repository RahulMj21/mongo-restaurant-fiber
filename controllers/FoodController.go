package controllers

import (
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var FoodCollection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFood(c *fiber.Ctx) error {
	idParam := c.Params("id")

	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "food not found",
		})
	}

	food := models.Food{}

	err = FoodCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&food)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "food not found",
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   food,
	})
}

func CreateFood(c *fiber.Ctx) error {
	food := models.Food{}
	menu := models.Menu{}

	if err := c.BodyParser(&food); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	validationErr := validate.Struct(food)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": validationErr.Error(),
		})
	}

	menuObjectId, err := primitive.ObjectIDFromHex(food.MenuId)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	err = MenuCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: menuObjectId}}).Decode(&menu)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	price := ToFixed(food.Price, 2)
	food.Price = price

	insertedItem, err := FoodCollection.InsertOne(c.Context(), food)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	newFoodItem := models.Food{}
	err = FoodCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newFoodItem)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   newFoodItem,
	})

}

func ToFixed(num float64, precision int) float64 {
	return num
}
