package controllers

import (
	"strconv"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UsersCollection = database.OpenCollection(database.Client, "User")

func GetUsers(c *fiber.Ctx) error {
	recordPerPage, recordPerPageErr := strconv.Atoi(c.Query("recordPerPage"))
	if recordPerPageErr != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, pageErr := strconv.Atoi(c.Query("page"))
	if pageErr != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "first_name", Value: 1},
		{Key: "last_name", Value: 1},
		{Key: "email", Value: 1},
		{Key: "phone", Value: 1},
		{Key: "avatar", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$ROOT"}}},
	}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "users", Value: bson.D{{Key: "$slice", Value: []interface{}{startIndex, recordPerPage}}}},
	}}}

	cursor, err := UsersCollection.Aggregate(c.Context(), mongo.Pipeline{
		matchStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "failed to get users"})
	}

	users := []primitive.M{}
	if err := cursor.All(c.Context(), &users); err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "failed to get"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": users})
}

func GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "user_id cannot be empty"})
	}

	userId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "invalid user_id"})
	}

	user := models.User{}
	opt := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}, {Key: "access_token", Value: 0},
		{Key: "refresh_token", Value: 0},
	})

	err = UsersCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: userId}}, opt).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{"status": "fail", "message": "failed to fetch"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": user})
}

func Signup(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func Signin(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
