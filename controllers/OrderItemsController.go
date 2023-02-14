package controllers

import (
	"time"

	"github.com/RahulMj21/mongo-restaurant-fiber/database"
	"github.com/RahulMj21/mongo-restaurant-fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	allOrderItems, err := ItemsByOrderId(orderIdParam, c.Context())
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": allOrderItems})
}

func ItemsByOrderId(id string, ctx *fasthttp.RequestCtx) (OrderItems []primitive.M, err error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$toObjectId", Value: "$food_id"}}}}},
		{Key: "as", Value: "food"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$food"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "order"},
		{Key: "localField", Value: "order_id"},
		{Key: "foreignField", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$toObjectId", Value: "$order_id"}}}}},
		{Key: "as", Value: "food"},
	}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$order"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "table"},
		{Key: "localField", Value: "order.table_id"},
		{Key: "foreignField", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$toObjectId", Value: "$order.table_id"}}}}},
		{Key: "as", Value: "table"},
	}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$table"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "amount", Value: "$food.price"},
		{Key: "total_count", Value: 1},
		{Key: "food_name", Value: "$food.name"},
		{Key: "food_image", Value: "$food.image"},
		{Key: "table_number", Value: "$table.table_number"},
		{Key: "price", Value: "$food.price"},
		{Key: "quantity", Value: 1},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{
		{Key: "order_id", Value: "$order_id"},
		{Key: "table_id", Value: "$table_id"},
		{Key: "table_number", Value: "$table_number"},
	}},
		{Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
	}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "payment_due", Value: 1},
		{Key: "total_count", Value: 1},
		{Key: "table_number", Value: "$_id.table_number"},
		{Key: "order_items", Value: 1},
	}}}

	orderItems := []primitive.M{}
	cursor, err := OrderItemsCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		return orderItems, err
	}

	if err := cursor.All(ctx, &orderItems); err != nil {
		return orderItems, err
	}

	return orderItems, nil
}

func CreateOrderItem(c *fiber.Ctx) error {
	orderItemPack := OrderItemPack{}
	order := models.Order{}

	err := c.BodyParser(&orderItemPack)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}
	order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.TableId = &orderItemPack.TableId
	orderId, err := OrderItemOrderCreater(order)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	var orderItemsToBeInserted []interface{}

	for _, orderItem := range orderItemPack.OrderItems {
		orderItem.OrderId = orderId
		validationErr := Validate.Struct(orderItem)
		if validationErr != nil {
			return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": validationErr.Error()})
		}

		orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		unitPrice := toFixed(*orderItem.UnitPrice, 2)
		orderItem.UnitPrice = &unitPrice

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	insertedItems, err := OrderItemsCollection.InsertMany(c.Context(), orderItemsToBeInserted)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "failed to insert order items"})
	}

	return c.Status(201).JSON(&fiber.Map{"status": "success", "data": insertedItems})
}

func UpdateOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	orderItem := models.OrderItem{}
	if idParam == "" {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": "id cannot be empty"})
	}
	orderItemId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	if err := c.BodyParser(&orderItem); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	orderItemObj := primitive.D{}

	if orderItem.UnitPrice != nil {
		price := toFixed(*orderItem.UnitPrice, 2)
		orderItemObj = append(orderItemObj, bson.E{Key: "unit_price", Value: &price})
	}
	if orderItem.Quantity != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}
	if orderItem.FoodId != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "food_id", Value: orderItem.FoodId})
	}

	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItemObj = append(orderItemObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

	filter := bson.D{{Key: "_id", Value: orderItemId}}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := OrderItemsCollection.UpdateOne(c.Context(), filter, bson.D{{Key: "$set", Value: orderItemObj}}, &opt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": "cannot update order_item"})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": result})
}
