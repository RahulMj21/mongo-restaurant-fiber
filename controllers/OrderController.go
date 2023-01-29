package controllers

import "github.com/RahulMj21/mongo-restaurant-fiber/database"

var OrderCollection = database.OpenCollection(database.Client, "order")
