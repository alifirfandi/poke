package entity

import (
	"time"
)

type FightHistory struct {
	ID                 uint                 `json:"id" gorm:"primarykey"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	FightHistoryDetail []FightHistoryDetail `json:"fight_history_detail" gorm:"foreignKey:FightHistoryID"`
}
