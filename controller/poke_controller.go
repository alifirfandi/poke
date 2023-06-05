package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"math"
	"net/http"
	"pokeapi/helper"
	"pokeapi/middleware"
	"pokeapi/model"
	"pokeapi/service"
	"strconv"
	"time"
)

type PokeController struct {
	PokeService service.PokeService
}

func NewPokeController(pokeService *service.PokeService) PokeController {
	return PokeController{
		PokeService: *pokeService,
	}
}

func (c PokeController) Route(app fiber.Router) {
	app.Post("/login", c.Login)
	app.Get("/pokemon", middleware.AuthMiddleware, c.GetAll)
	app.Get("/pokemon/:name", middleware.AuthMiddleware, c.GetOne)
	app.Post("/fight", middleware.AuthMiddleware, c.Fight)
	app.Get("/fight/history", middleware.AuthMiddleware, c.GetHistories)
	app.Put("/cancel", middleware.AuthMiddleware, c.CancelPokemon)
	app.Get("/leaderboard", middleware.AuthMiddleware, c.Leaderboard)
}

func (c PokeController) GetAll(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "operational" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	page := ctx.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	pokeData, pokeApiRes, err := c.PokeService.GetListPokemon(pageInt)

	return ctx.Status(http.StatusOK).JSON(model.ResponseList{
		Page:      pageInt,
		PageTotal: int(math.Ceil(float64(pokeApiRes.Count) / 10.0)),
		Data:      pokeData,
		DataTotal: pokeApiRes.Count,
	})
}

func (c PokeController) GetOne(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "operational" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	name := ctx.Params("name")
	pokeData, err := c.PokeService.GetPokemonData(name)
	if err != nil {
		return ctx.Status(404).JSON(model.Response{
			Error: "Data Not Found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: pokeData,
	})
}

func (c PokeController) Fight(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	var reqBody model.PokemonCreateReqBody
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(model.Response{
			Error: "Bad Request",
		})
	}

	if helper.HasDuplicateString(reqBody.Pokemon) {
		return ctx.Status(400).JSON(model.Response{
			Error: "Bad Request",
		})
	}

	pokeData, err := c.PokeService.FightPokemon(reqBody.Pokemon)
	if err != nil {
		return ctx.Status(400).JSON(model.Response{
			Error: "Bad Request",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: pokeData,
	})
}

func (c PokeController) GetHistories(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "superadmin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	var req model.PokemonReqQuery
	if ctx.Query("start_date") != "" && ctx.Query("end_date") != "" {
		req = model.PokemonReqQuery{
			StartDate: fmt.Sprintf("%s 00:00:00", ctx.Query("start_date")),
			EndDate:   fmt.Sprintf("%s 23:59:59", ctx.Query("end_date")),
		}
	}

	fightHistories, err := c.PokeService.FightHistories(req)
	if err != nil {
		return ctx.Status(500).JSON(model.Response{
			Error: "Internal Server Error",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: fightHistories,
	})
}

func (c PokeController) Leaderboard(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "superadmin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	leaderboardData, err := c.PokeService.GetLeaderboard()
	if err != nil {
		return ctx.Status(500).JSON(model.Response{
			Error: "Internal Server Error",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: leaderboardData,
	})
}

func (c PokeController) CancelPokemon(ctx *fiber.Ctx) error {
	role := ctx.Locals("role")
	if role != "superadmin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	var reqBody model.PokemonCancelReqBody
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(model.Response{
			Error: "Bad Request",
		})
	}

	pokeData, err := c.PokeService.CancelPokemon(reqBody)
	if err != nil {
		return ctx.Status(400).JSON(model.Response{
			Error: "Bad Request",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: pokeData,
	})
}

func (c PokeController) Login(ctx *fiber.Ctx) error {
	user := new(model.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if user.Role != "operational" && user.Role != "admin" && user.Role != "superadmin" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{
		"token": tokenString,
	})
}
