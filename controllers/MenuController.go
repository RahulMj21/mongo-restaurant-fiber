package controllers

import "github.com/RahulMj21/mongo-restaurant-fiber/database"

var MenuCollection = database.OpenCollection(database.Client, "menu")
