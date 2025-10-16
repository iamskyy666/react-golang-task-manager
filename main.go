package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct{
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

func main() {
	fmt.Println("Hello World!!")
	app:=fiber.New()

	err:= godotenv.Load(".env")
	if err!=nil{
		log.Fatal("ERROR loading .env file:",err)
	}

	PORT:=os.Getenv("PORT")

	todos:=[]Todo{}



	// routes
	// get list of TODOs
	app.Get("/api/todos",func (c *fiber.Ctx) error {
		// return c.Status(200).JSON(fiber.Map{"msg":"Testing GET route.."})
		 return c.Status(200).JSON(todos)
	})

	// create TODO
	app.Post("/api/todos",func (c *fiber.Ctx)error  {
		todo:= &Todo{} // empty values for now
		if err:=c.BodyParser(todo); err !=nil{
			return err
		}

		if todo.Body == ""{
			return c.Status(400).JSON(fiber.Map{"error":"Todo - body is required!"})
		}

		todo.ID = len(todos)+1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// update a TODO
	app.Patch("/api/todos/:id",func (c *fiber.Ctx)error{
		id := c.Params("id")

		for i, todo := range todos{
			if fmt.Sprint(todo.ID) == id{
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error":"Todo not found!"})
	})

	// Delete a TODO
	app.Delete("/api/todos/:id", func (c *fiber.Ctx)error  {
		id := c.Params("id")
		for i, todo := range todos{
			if fmt.Sprint(todo.ID) == id{
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(201).JSON(fiber.Map{"error":"Todo deleted ✅"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error":"Todo not found!"})
		
	})



	log.Fatal(app.Listen(":"+PORT))
}
