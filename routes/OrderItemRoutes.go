package routes

import (
	"github.com/RahulMj21/mongo-restaurant-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func OrderItemRoutes(api fiber.Router) {
	api.Get("/order-items", controllers.GetOrderItems)
	api.Get("/order-items/:id", controllers.GetOrderItem)
	api.Get("/order-items-by-order/:order_id", controllers.GetOrderItemsByOrderId)
	api.Post("/order-items", controllers.CreateOrderItem)
	api.Put("/order-items", controllers.UpdateOrderItem)
}
