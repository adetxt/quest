package questmysql

import (
	"time"

	"github.com/adetxt/quest/domain"
)

type Quest struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

type Objective struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	QuestID     int64  `gorm:"column:quest_id"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

type QuestUser struct {
	UserID    int64     `gorm:"column:user_id;uniqueIndex:idx_user_id_quest_id"`
	QuestID   int64     `gorm:"column:quest_id;uniqueIndex:idx_user_id_quest_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type ObjectiveUser struct {
	UserID      int64     `gorm:"column:user_id;uniqueIndex:idx_user_id_objective_id"`
	ObjectiveID int64     `gorm:"column:objective_id;uniqueIndex:idx_user_id_objective_id"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (i *Quest) ToEntity() *domain.Quest {
	return &domain.Quest{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
	}
}

func (i *Objective) ToEntity() *domain.Objective {
	return &domain.Objective{
		ID:          i.ID,
		QuestID:     i.QuestID,
		Title:       i.Title,
		Description: i.Description,
	}
}
