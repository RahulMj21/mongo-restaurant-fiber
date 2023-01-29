package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func InvoiceRoutes(api fiber.Router) {
	api.Get("/invoices", controllers.GetInvoices)
	api.Get("/invoices/:id", controllers.GetInvoice)
	api.Post("/invoices", controllers.CreateInvoice)
	api.Patch("/invoices/:id", controllers.UpdateInvoice)

}
