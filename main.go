package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        int    `json:"_id" bson:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello World")
	err:= godotenv.Load(".env")
	if err!=nil{
		log.Fatal("⚠️ ERR. loading the .env file:",err)
	}

	MONGODB_URI:=os.Getenv("MONGODB_URI")
	clientOptions:=options.Client().ApplyURI(MONGODB_URI) // connect to MongoDB

	client,err:= mongo.Connect(context.Background(),clientOptions)
	if err!=nil{
		log.Fatal("⚠️ ERR:",err)
	}
	err = client.Ping(context.Background(),nil)
	if err!=nil{
		log.Fatal("⚠️ ERR:",err)
	}
	fmt.Println("Connected to MongoDB Atlas.. ✅")

	collection = client.Database("golang_db").Collection("todos")

	app:= fiber.New()
	app.Get("/api/todos", GetTodos)
	app.Post("/api/todos",CreateTodo)
	app.Patch("/api/todos/:id",UpdateTodo)
	app.Post("/api/todos/:id",DeleteTodo)

	PORT:= os.Getenv("PORT")
	if PORT == ""{
		PORT="5000"
	}

	log.Fatal(app.Listen("0.0.0.0:"+PORT))
	
}