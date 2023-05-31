package entity

type FightHistoryDetail struct {
	ID             uint   `json:"id" gorm:"primarykey"`
	FightHistoryID uint   `json:"id_fight_history" gorm:"foreignKey:FightHistoryID"`
	Pokemon        string `json:"pokemon"`
	Score          int    `json:"score"`
}
