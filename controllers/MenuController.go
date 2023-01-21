package controllers

import (
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MenuCollection = database.OpenCollection(database.Client, "menu")
var Validate = validator.New()

func GetMenus(c *fiber.Ctx) error {
	var menus []models.Menu

	cursor, err := MenuCollection.Find(c.Context(), bson.D{{}})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := cursor.All(c.Context(), &menus); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   menus,
	})
}

func GetMenu(c *fiber.Ctx) error {
	menuIdParam := c.Params("id")

	menuId, err := primitive.ObjectIDFromHex(menuIdParam)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	menu := models.Menu{}
	err = MenuCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: menuId}}).Decode(&menu)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   menu,
	})
}

func CreateMenu(c *fiber.Ctx) error {
	menu := models.Menu{}
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	validationErr := Validate.Struct(&menu)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": validationErr.Error(),
		})
	}

	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertedItem, err := MenuCollection.InsertOne(c.Context(), menu)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	newItem := models.Menu{}
	err = MenuCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newItem)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(&fiber.Map{
		"status": "success",
		"data":   newItem,
	})
}

func UpdateMenu(c *fiber.Ctx) error {
	menuIdParam := c.Params("id")
	menu := models.Menu{}

	menuId, err := primitive.ObjectIDFromHex(menuIdParam)
	if err != nil {
		return c.Status(200).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	err = c.BodyParser(&menu)
	if err != nil {
		return c.Status(200).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var menuObj primitive.D

	if menu.StartDate != nil && menu.EndDate != nil {
		if !inTimeSpan(*menu.StartDate, *menu.EndDate) {
			return c.Status(400).JSON(&fiber.Map{
				"status":  "fail",
				"message": "please retype the time",
			})
		}

		menuObj = append(menuObj, bson.E{Key: "start_date", Value: menu.StartDate})
		menuObj = append(menuObj, bson.E{Key: "end_date", Value: menu.EndDate})

		if menu.Name != "" {
			menuObj = append(menuObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			menuObj = append(menuObj, bson.E{Key: "category", Value: menu.Category})
		}

		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		menuObj = append(menuObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := MenuCollection.UpdateOne(c.Context(), bson.D{{Key: "_id", Value: menuId}}, menuObj, &opt)
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
	return c.Status(500).JSON(&fiber.Map{
		"status":  "fail",
		"message": "update failed",
	})
}

func inTimeSpan(start, end time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
