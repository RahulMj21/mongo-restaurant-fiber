package controllers

import "github.com/RahulMj21/mongo-restaurant-fiber/database"

var TableCollection = database.OpenCollection(database.Client, "table")
