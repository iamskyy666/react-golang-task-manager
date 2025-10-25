package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
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

	defer client.Disconnect(context.Background()) // optimization

	err = client.Ping(context.Background(),nil)
	if err!=nil{
		log.Fatal("⚠️ ERR:",err)
	}
	fmt.Println("Connected to MongoDB Atlas.. ✅")

	collection = client.Database("golang_db").Collection("todos")

	app:= fiber.New()
	//! CRUD endpoints
	app.Get("/api/todos", GetTodos)
	app.Post("/api/todos",CreateTodo)
	app.Patch("/api/todos/:id",UpdateTodo)
	app.Delete("/api/todos/:id",DeleteTodo)

	PORT:= os.Getenv("PORT")
	if PORT == ""{
		PORT="5000"
	}

	log.Fatal(app.Listen("0.0.0.0:"+PORT))
	
}

// Add fiber for CRUD Ops.

//! READ
func GetTodos(ctx *fiber.Ctx)error{
var todos []Todo
cursor,err:= collection.Find(context.Background(),bson.M{}) // No filters - Fetching all todos/docs
if err!=nil{
	return  err
}

defer cursor.Close(context.Background()) // optimization

// If no error, then..
for cursor.Next(context.Background()){
	var todo Todo
	if err:= cursor.Decode(&todo); err!=nil{
		return  err
	}
	todos = append(todos, todo)
}
return ctx.JSON(todos)
}

//! CREATE
func CreateTodo(ctx *fiber.Ctx)error{
	todo:= new(Todo)
	if err:= ctx.BodyParser(todo); err != nil{
		return err
	}
	// validate
	if todo.Body==""{
		return ctx.Status(400).JSON(fiber.Map{"⚠️ERROR":"Todo-body cannot be empty!"})
	}

	// insert into MongoDB
	insertedResult, err:=collection.InsertOne(context.Background(),todo)
	if err!=nil{
		return err
	}
	todo.ID = insertedResult.InsertedID.(primitive.ObjectID)
	return ctx.Status(201).JSON(todo)
}


//! UPDATE
 func UpdateTodo(ctx *fiber.Ctx)error{
	id:= ctx.Params("id")
	// ObjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a valid ObjectID.
	objId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		return ctx.Status(400).JSON(fiber.Map{"⚠️ERROR":"Invalid todo-ID!"})
	}
	filter:=bson.M{"_id":objId}
	update:= bson.M{"$set":bson.M{"completed":true}} // change/update todo status
	_,err=collection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{"success":"true"})

 }	

//! DELETE
	func DeleteTodo(ctx *fiber.Ctx)error{
		id:= ctx.Params("id")

	objId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		return ctx.Status(400).JSON(fiber.Map{"⚠️ERROR":"Invalid todo-ID!"})
	}

	filter:= bson.M{"_id":objId}

	// If no-error, then delete 1 from the collection
	_,err=collection.DeleteOne(context.Background(),filter)
	if err!=nil{
		return err
	}

	// if no-error while deleting, then, return a success-response
	return ctx.Status(200).JSON(fiber.Map{"success":"true"})
}

// ⏱️
