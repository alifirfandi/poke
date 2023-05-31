package pokemon

import (
	"math"
	"pokeapi/model"
)

type Pokemon struct{}

func New() *Pokemon {
	return &Pokemon{}
}

func (p Pokemon) PokemonDataSourceToListString(pokeDataSource model.PokeDataSourceRes) []string {
	var listPokemon []string
	for _, p := range pokeDataSource.Results {
		listPokemon = append(listPokemon, p.Name)
	}
	return listPokemon
}

func (p Pokemon) PokemonDetailDataSourceToPokemon(pokeDataSource model.PokeDetailDataSourceRes) model.Pokemon {
	pokemon := model.Pokemon{
		Name: pokeDataSource.Name,
	}
	var cp float64
	for _, p := range pokeDataSource.Stats {
		cp += float64(p.BaseStat)
		pokemon.Stats = append(pokemon.Stats, model.Stat{
			Name:  p.Stat.Name,
			Value: p.BaseStat,
		})
	}
	pokemon.CombatPower = math.Round(cp/float64(len(pokeDataSource.Stats))*100) / 100

	return pokemon
}

func (p Pokemon) FightPokemon(pokemons []model.Pokemon) []model.Pokemon {
	n := len(pokemons)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if pokemons[j].CombatPower < pokemons[j+1].CombatPower {
				pokemons[j], pokemons[j+1] = pokemons[j+1], pokemons[j]
			}
		}
	}

	return pokemons
}
