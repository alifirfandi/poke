package pokemon_test

import (
	"github.com/stretchr/testify/assert"
	"pokeapi/model"
	"pokeapi/pokemon"
	"testing"
)

func TestPokemonDataSourceToListString(t *testing.T) {
	p := pokemon.New()
	testTable := []struct {
		pokeDataSource  model.PokeDataSourceRes
		expectedOutcome []string
	}{
		{
			pokeDataSource:  model.PokeDataSourceRes{},
			expectedOutcome: nil,
		},
		{
			pokeDataSource: model.PokeDataSourceRes{
				Count: 5,
				Results: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					{Name: "bulbasaur", URL: "https://pokeapi.co/api/v2/pokemon/1/"},
					{Name: "ivysaur", URL: "https://pokeapi.co/api/v2/pokemon/2/"},
					{Name: "venusaur", URL: "https://pokeapi.co/api/v2/pokemon/3/"},
					{Name: "charmander", URL: "https://pokeapi.co/api/v2/pokemon/4/"},
					{Name: "charmeleon", URL: "https://pokeapi.co/api/v2/pokemon/5/"},
				},
			},
			expectedOutcome: []string{"bulbasaur", "ivysaur", "venusaur", "charmander", "charmeleon"},
		},
	}

	for _, test := range testTable {
		result := p.PokemonDataSourceToListString(test.pokeDataSource)
		assert.Equal(t, test.expectedOutcome, result)
	}
}

func TestPokemonDetailDataSourceToPokemon(t *testing.T) {
	p := pokemon.New()
	testTable := []struct {
		pokeDataSource  model.PokeDetailDataSourceRes
		expectedOutcome model.Pokemon
	}{
		{
			pokeDataSource: model.PokeDetailDataSourceRes{
				Name: "bulbasaur",
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{
						BaseStat: 45,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "hp"},
					},
					{
						BaseStat: 49,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "attack"},
					},
					{
						BaseStat: 49,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "defense"},
					},
				},
			},
			expectedOutcome: model.Pokemon{
				Name: "bulbasaur",
				Stats: []model.Stat{
					{
						Name:  "hp",
						Value: 45,
					}, {
						Name:  "attack",
						Value: 49,
					}, {
						Name:  "defense",
						Value: 49,
					},
				},
				CombatPower: 47.67,
			},
		},
		{
			pokeDataSource: model.PokeDetailDataSourceRes{
				Name: "ivysaur",
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{
						BaseStat: 60,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "hp"},
					},
					{
						BaseStat: 62,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "attack"},
					},
					{
						BaseStat: 63,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "defense"},
					},
				},
			},
			expectedOutcome: model.Pokemon{
				Name: "ivysaur",
				Stats: []model.Stat{
					{
						Name:  "hp",
						Value: 60,
					}, {
						Name:  "attack",
						Value: 62,
					}, {
						Name:  "defense",
						Value: 63,
					},
				},
				CombatPower: 61.67,
			},
		},
	}

	for _, test := range testTable {
		result := p.PokemonDetailDataSourceToPokemon(test.pokeDataSource)
		assert.Equal(t, test.expectedOutcome, result)
	}
}

func TestFightPokemon(t *testing.T) {
	p := pokemon.New()
	testTable := []struct {
		pokemons        []model.Pokemon
		expectedOutcome []model.Pokemon
	}{
		{
			pokemons: []model.Pokemon{
				{
					Name:        "bulbasaur",
					Stats:       nil,
					CombatPower: 100,
				}, {
					Name:        "pikachu",
					Stats:       nil,
					CombatPower: 60,
				}, {
					Name:        "ivysaur",
					Stats:       nil,
					CombatPower: 70,
				}, {
					Name:        "charizard",
					Stats:       nil,
					CombatPower: 20,
				}, {
					Name:        "snorlax",
					Stats:       nil,
					CombatPower: 150,
				},
			},
			expectedOutcome: []model.Pokemon{
				{
					Name:        "snorlax",
					Stats:       nil,
					CombatPower: 150,
				}, {
					Name:        "bulbasaur",
					Stats:       nil,
					CombatPower: 100,
				}, {
					Name:        "ivysaur",
					Stats:       nil,
					CombatPower: 70,
				}, {
					Name:        "pikachu",
					Stats:       nil,
					CombatPower: 60,
				}, {
					Name:        "charizard",
					Stats:       nil,
					CombatPower: 20,
				},
			},
		},
		{
			pokemons: []model.Pokemon{
				{
					Name:        "bulbasaur",
					Stats:       nil,
					CombatPower: 30,
				}, {
					Name:        "pikachu",
					Stats:       nil,
					CombatPower: 50,
				}, {
					Name:        "ivysaur",
					Stats:       nil,
					CombatPower: 20,
				}, {
					Name:        "charizard",
					Stats:       nil,
					CombatPower: 70,
				}, {
					Name:        "snorlax",
					Stats:       nil,
					CombatPower: 150,
				},
			},
			expectedOutcome: []model.Pokemon{
				{
					Name:        "snorlax",
					Stats:       nil,
					CombatPower: 150,
				}, {
					Name:        "charizard",
					Stats:       nil,
					CombatPower: 70,
				}, {
					Name:        "pikachu",
					Stats:       nil,
					CombatPower: 50,
				}, {
					Name:        "bulbasaur",
					Stats:       nil,
					CombatPower: 30,
				}, {
					Name:        "ivysaur",
					Stats:       nil,
					CombatPower: 20,
				},
			},
		},
	}

	for _, test := range testTable {
		result := p.FightPokemon(test.pokemons)
		assert.Equal(t, test.expectedOutcome, result)
	}
}
