package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"pokeapi/entity"
	"pokeapi/model"
	"time"
)

type PokeRepository struct {
	DB *gorm.DB
}

func NewPokeRepository(mysql *gorm.DB) PokeRepository {
	return PokeRepository{
		DB: mysql,
	}
}

func (r PokeRepository) GetAllPokemon(offset int) (model.PokeDataSourceRes, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?limit=10&offset=%d", offset)

	response, err := http.Get(url)
	if err != nil {
		return model.PokeDataSourceRes{}, err
	}
	defer response.Body.Close()

	var pokeApiRes model.PokeDataSourceRes
	err = json.NewDecoder(response.Body).Decode(&pokeApiRes)
	if err != nil {
		return model.PokeDataSourceRes{}, err
	}

	return pokeApiRes, nil
}

func (r PokeRepository) GetOnePokemon(name string) (model.PokeDetailDataSourceRes, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	response, err := http.Get(url)
	if err != nil {
		return model.PokeDetailDataSourceRes{}, err
	}
	defer response.Body.Close()

	var pokeApi model.PokeDetailDataSourceRes
	err = json.NewDecoder(response.Body).Decode(&pokeApi)
	if err != nil {
		return model.PokeDetailDataSourceRes{}, err
	}

	return pokeApi, nil
}

func (r PokeRepository) InsertFightHistory() (entity.FightHistory, error) {
	var fightHistory entity.FightHistory
	res := r.DB.Create(&fightHistory)
	if res.RowsAffected == 0 {
		return entity.FightHistory{}, errors.New("failed insert fight history data")
	}
	return fightHistory, nil
}

func (r PokeRepository) InsertFightHistoryDetail(fightHistoryDetail []entity.FightHistoryDetail) ([]entity.FightHistoryDetail, error) {
	res := r.DB.CreateInBatches(&fightHistoryDetail, len(fightHistoryDetail))
	if res.RowsAffected < int64(len(fightHistoryDetail)) {
		return []entity.FightHistoryDetail{}, errors.New("failed insert fight history detail data in batch")
	}
	return fightHistoryDetail, nil
}

func (r PokeRepository) GetFightHistory(req model.PokemonReqQuery) ([]entity.FightHistory, error) {
	var fightHistories []entity.FightHistory
	db := r.DB.Preload("FightHistoryDetail")

	if req.StartDate != "" && req.EndDate != "" {
		_, err := time.Parse("2006-01-02 15:04:05", req.StartDate)
		if err != nil {
			return []entity.FightHistory{}, err
		}

		_, err = time.Parse("2006-01-02 15:04:05", req.EndDate)
		if err != nil {
			return []entity.FightHistory{}, err
		}

		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	err := db.Order("id DESC").Find(&fightHistories).Error
	if err != nil {
		return []entity.FightHistory{}, err
	}

	return fightHistories, nil
}

func (r PokeRepository) GetSumScore() ([]model.Leaderboard, error) {
	var leaderboard []model.Leaderboard
	_ = r.DB.Table("fight_history_details").
		Select("pokemon, SUM(score) as total_score").
		Group("pokemon").
		Order("total_score DESC").
		Scan(&leaderboard)

	return leaderboard, nil
}

func (r PokeRepository) CancelScorePokemon(req model.PokemonCancelReqBody) (entity.FightHistoryDetail, error) {
	var fightHistoryDetail entity.FightHistoryDetail
	err := r.DB.Where("fight_history_id = ? AND pokemon = ? AND score != ?", req.FightHistoryID, req.Pokemon, 0).
		First(&fightHistoryDetail).Error
	if err != nil {
		return entity.FightHistoryDetail{}, err
	}

	err = r.DB.Model(&entity.FightHistoryDetail{}).
		Where("fight_history_id = ? AND score < ? AND score > ?", req.FightHistoryID, fightHistoryDetail.Score, 0).
		Update("score", gorm.Expr("score + ?", 1)).Error
	if err != nil {
		return entity.FightHistoryDetail{}, err
	}

	fightHistoryDetail.Score = 0
	err = r.DB.Save(&fightHistoryDetail).Error
	if err != nil {
		return entity.FightHistoryDetail{}, err
	}

	return fightHistoryDetail, nil
}
