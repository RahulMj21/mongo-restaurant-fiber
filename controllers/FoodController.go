package controllers

import (
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
	// food := models.Food{}
	// menu := models.Menu{}

	return c.SendStatus(200)

}
