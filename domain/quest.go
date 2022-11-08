package domain

import "context"

type QuestUsecase interface {
	GetUserQuests(ctx context.Context, params *GetQuestsParam) ([]*QuestUser, error)
	AssignQuest(ctx context.Context, questID, userID int64) error
	RevokeQuest(ctx context.Context, questID, userID int64) error
	CompleteObjective(ctx context.Context, objectiveID, userID int64) error
}

type QuestRepository interface {
	GetUserQuets(ctx context.Context, params *GetQuestsParam) ([]*QuestUser, error)
	AssignQuest(ctx context.Context, questID, userID int64) error
	RevokeQuest(ctx context.Context, questID, userID int64) error
	SetObjectiveStatus(ctx context.Context, objectiveID, userID int64, status ObjectiveStatus) error
}

type ObjectiveStatus string

const (
	ObjectiveInprogress ObjectiveStatus = "IN_PROGRESS"
	ObjectiveComplete   ObjectiveStatus = "COMPLETE"
)

func (v ObjectiveStatus) String() string {
	return string(v)
}

type Quest struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Objective struct {
	ID          int64  `json:"id"`
	QuestID     int64  `json:"quest_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QuestUser struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Objectives  []ObjectiveUser `json:"objectives"`
	Percentage  float32         `json:"percentage"`
}

type ObjectiveUser struct {
	ID          int64           `json:"id"`
	QuestID     int64           `json:"quest_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      ObjectiveStatus `json:"status"`
}

type GetQuestsParam struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}
