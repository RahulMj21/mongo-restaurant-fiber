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

type InvoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

var InvoiceCollection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(c *fiber.Ctx) error {
	cursor, err := InvoiceCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var invoices []primitive.M

	if err := cursor.All(c.Context(), &invoices); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   invoices,
	})
}

func GetInvoice(c *fiber.Ctx) error {
	idParam := c.Params("id")

	invoiceId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "invoice not found",
		})
	}

	invoice := models.Invoice{}
	err = InvoiceCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: invoiceId}}).Decode(&invoice)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "invoice not found",
		})
	}

	invoiceView := InvoiceViewFormat{}

	allOrderItems, err := ItemsByOrderId(idParam, c.Context())
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	invoiceView.InvoiceId = invoice.ID.Hex()
	invoiceView.OrderId = invoice.OrderId
	invoiceView.PaymentDueDate = invoice.PaymentDueDate
	invoiceView.PaymentStatus = invoice.PaymentStatus

	invoiceView.PaymentMethod = "null"
	if invoice.PaymentMethod != nil {
		invoiceView.PaymentMethod = *invoice.PaymentMethod
	}

	invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
	invoiceView.TableNumber = allOrderItems[0]["table_number"]
	invoiceView.OrderDetails = allOrderItems[0]["order_items"]

	return c.Status(200).JSON(&fiber.Map{
		"status": "success",
		"data":   invoiceView,
	})
}

func CreateInvoice(c *fiber.Ctx) error {
	invoice := models.Invoice{}
	order := models.Order{}

	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	invoiceId, err := primitive.ObjectIDFromHex(invoice.OrderId)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	if err := OrderCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: invoiceId}}).Decode(&order); err != nil {
		return c.Status(404).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	invoice.ID = primitive.NewObjectID()
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))

	status := "PENDING"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &status
	}

	validationErr := Validate.Struct(invoice)
	if validationErr != nil {
		return c.Status(400).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	insertedItem, err := InvoiceCollection.InsertOne(c.Context(), invoice)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	newInvoice := models.Invoice{}

	if err := InvoiceCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newInvoice); err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(201).JSON(&fiber.Map{"status": "success", "data": newInvoice})
}

func UpdateInvoice(c *fiber.Ctx) error {
	idParam := c.Params("id")

	invoiceId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	invoice := models.Invoice{}
	order := models.Order{}

	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	orderId, err := primitive.ObjectIDFromHex(invoice.OrderId)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "order not found for the invoice: " + idParam,
		})
	}

	if err := OrderCollection.FindOne(c.Context(), bson.D{{Key: "_id", Value: orderId}}).Decode(&order); err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"status":  "fail",
			"message": "order not found for the invoice: " + idParam,
		})
	}

	filter := bson.D{{Key: "_id", Value: invoiceId}}
	var invoiceObj primitive.D

	if invoice.PaymentMethod != nil {
		invoiceObj = append(invoiceObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
	}
	if invoice.PaymentStatus != nil {
		invoiceObj = append(invoiceObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
	}

	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoiceObj = append(invoiceObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "Pending"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &status
	}

	result, err := InvoiceCollection.UpdateOne(c.Context(), filter, bson.D{{Key: "$set", Value: invoiceObj}}, &opt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(200).JSON(&fiber.Map{"status": "success", "data": result})
}
