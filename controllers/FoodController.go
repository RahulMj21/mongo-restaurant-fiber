package controllers

import (
	"math"
	"strconv"
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FoodCollection = database.OpenCollection(database.Client, "food")

func GetFoods(c *fiber.Ctx) error {
	resultPerPage, err := strconv.Atoi(c.Query("resultPerPage"))
	if err != nil || resultPerPage < 1 {
		resultPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * resultPerPage
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "null"},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: "1"}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
	}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "food_items", Value: bson.D{
			{Key: "$slice", Value: []interface{}{"$data", startIndex, resultPerPage}},
		}},
	}}}

	result, err := FoodCollection.Aggregate(c.Context(), mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var foods []primitive.D

	if err := result.All(c.Context(), &foods); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   foods,
	})
}

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

	validationErr := Validate.Struct(food)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": validationErr.Error(),
		})
	}

	menuIdHex := food.MenuId
	menuObjectId, err := primitive.ObjectIDFromHex(*menuIdHex)
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
	price := toFixed(*food.Price, 2)
	food.Price = &price

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

func UpdateFood(c *fiber.Ctx) error {
	foodIdParam := c.Params("id")
	foodId, err := primitive.ObjectIDFromHex(foodIdParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	food := models.Food{}
	menu := models.Menu{}

	if err := c.BodyParser(&food); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var foodObj primitive.D

	if food.Name != nil {
		foodObj = append(foodObj, bson.E{Key: "name", Value: food.Name})
	}
	if food.Price != nil {
		foodObj = append(foodObj, bson.E{Key: "price", Value: toFixed(*food.Price, 2)})
	}
	if food.FoodImage != nil {
		foodObj = append(foodObj, bson.E{Key: "food_image", Value: food.FoodImage})
	}

	if food.MenuId != nil {
		menuId, err := primitive.ObjectIDFromHex(*food.MenuId)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		err = MenuCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: menuId}}).Decode(&menu)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}

		foodObj = append(foodObj, bson.E{Key: "menu_id", Value: menu.ID.Hex()})
	}

	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	foodObj = append(foodObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	filter := bson.D{{Key: "_id", Value: foodId}}

	result, err := FoodCollection.UpdateOne(c.Context(), filter, bson.D{{Key: "$set", Value: foodObj}}, &opt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func round(num float64) int {
	return int(math.Round(num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(output)) / output
}
