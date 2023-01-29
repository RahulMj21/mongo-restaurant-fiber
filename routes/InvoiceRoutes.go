package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func InvoiceRoutes(app *fiber.App) {
	app.Get("/invoices", controllers.GetInvoices)
	app.Get("/invoices/:id", controllers.GetInvoice)
	app.Post("/invoices", controllers.CreateInvoice)
	app.Patch("/invoices/:id", controllers.UpdateInvoice)

}
