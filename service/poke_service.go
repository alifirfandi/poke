package service

import (
	"fmt"
	"pokeapi/entity"
	"pokeapi/model"
	"pokeapi/pokemon"
	"pokeapi/repository"
	"sync"
)

type PokeService struct {
	Pokemon        pokemon.Pokemon
	PokeRepository repository.PokeRepository
}

func NewPokeService(pokeRepository *repository.PokeRepository) PokeService {
	return PokeService{
		PokeRepository: *pokeRepository,
	}
}

func (s PokeService) GetListPokemon(page int) ([]string, model.PokeDataSourceRes, error) {
	offset := (page - 1) * 10
	pokeApiRes, err := s.PokeRepository.GetAllPokemon(offset)
	if err != nil {
		return nil, model.PokeDataSourceRes{}, err
	}
	pokeRes := s.Pokemon.PokemonDataSourceToListString(pokeApiRes)
	return pokeRes, pokeApiRes, nil
}

func (s PokeService) GetPokemonData(name string) (model.Pokemon, error) {
	pokeApiDetailRes, err := s.PokeRepository.GetOnePokemon(name)
	if err != nil {
		return model.Pokemon{}, err
	}
	pokeDetailRes := s.Pokemon.PokemonDetailDataSourceToPokemon(pokeApiDetailRes)
	return pokeDetailRes, nil
}

func (s PokeService) FightPokemon(pokemon []string) ([]model.Pokemon, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error
	var listPoke []model.Pokemon

	for _, p := range pokemon {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			pokeData, err := s.GetPokemonData(name)
			if err != nil {
				return
			}
			mutex.Lock()
			listPoke = append(listPoke, pokeData)
			mutex.Unlock()
		}(p)
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	if len(listPoke) != len(pokemon) {
		return nil, fmt.Errorf("failed to fetch all Pokemon data")
	}

	result := s.Pokemon.FightPokemon(listPoke)

	fightHistory, err := s.PokeRepository.InsertFightHistory()
	if err != nil {
		return nil, err
	}

	var detailFightData []entity.FightHistoryDetail
	score := 5
	for _, r := range result {
		detailFightData = append(detailFightData, entity.FightHistoryDetail{
			FightHistoryID: fightHistory.ID,
			Pokemon:        r.Name,
			Score:          score,
		})
		score--
	}
	_, err = s.PokeRepository.InsertFightHistoryDetail(detailFightData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s PokeService) FightHistories(req model.PokemonReqQuery) ([]entity.FightHistory, error) {
	fightHistories, err := s.PokeRepository.GetFightHistory(req)
	if err != nil {
		return []entity.FightHistory{}, err
	}

	return fightHistories, nil
}

func (s PokeService) GetLeaderboard() ([]model.Leaderboard, error) {
	leaderboardData, err := s.PokeRepository.GetSumScore()
	if err != nil {
		return []model.Leaderboard{}, err
	}

	return leaderboardData, nil
}

func (s PokeService) CancelPokemon(req model.PokemonCancelReqBody) (entity.FightHistoryDetail, error) {
	fightHistoryDetail, err := s.PokeRepository.CancelScorePokemon(req)
	if err != nil {
		return entity.FightHistoryDetail{}, err
	}

	return fightHistoryDetail, nil
}
