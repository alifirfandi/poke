package model

type Pokemon struct {
	Name        string  `json:"name"`
	Stats       []Stat  `json:"stats"`
	CombatPower float64 `json:"combat_power"`
}

type Stat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type PokeDataSourceRes struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokeDetailDataSourceRes struct {
	Name  string `json:"name"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

type PokemonReqQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type PokemonCreateReqBody struct {
	Pokemon []string `json:"pokemon"`
}

type PokemonCancelReqBody struct {
	FightHistoryID int    `json:"fight_history_id"`
	Pokemon        string `json:"pokemon"`
}

type Leaderboard struct {
	Pokemon    string `json:"pokemon"`
	TotalScore int    `json:"total_score"`
}
