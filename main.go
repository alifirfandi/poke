package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"os"
	"pokeapi/config"
	"pokeapi/controller"
	"pokeapi/repository"
	"pokeapi/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := config.Connect()
	if err != nil {
		panic(err)
	}

	pokeRepository := repository.NewPokeRepository(db)
	pokeService := service.NewPokeService(&pokeRepository)
	pokeController := controller.NewPokeController(&pokeService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(
		cors.Config{
			Next:             nil,
			AllowOrigins:     "*",
			AllowMethods:     "OPTIONS,GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           0,
		},
	))

	v1 := app.Group("/")
	pokeController.Route(v1)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf("%s:%s", host, port))
}
