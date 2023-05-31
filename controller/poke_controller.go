package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"math"
	"net/http"
	"pokeapi/helper"
	"pokeapi/model"
	"pokeapi/service"
	"strconv"
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
	app.Get("/pokemon", c.GetAll)
	app.Get("/pokemon/:name", c.GetOne)
	app.Post("/fight", c.Fight)
	app.Get("/fight/history", c.GetHistories)
	app.Put("/cancel", c.CancelPokemon)
	app.Get("/leaderboard", c.Leaderboard)
}

func (c PokeController) GetAll(ctx *fiber.Ctx) error {
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
	name := ctx.Params("name")
	pokeData, err := c.PokeService.GetPokemonData(name)
	if err != nil {
		return ctx.Status(404).JSON(model.Response{
			Error: "Data Pokemon Tidak Ditemukan",
		})
	}

	return ctx.Status(http.StatusOK).JSON(model.Response{
		Data: pokeData,
	})
}

func (c PokeController) Fight(ctx *fiber.Ctx) error {
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
