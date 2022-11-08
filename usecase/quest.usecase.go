package usecase

import (
	"context"

	"github.com/adetxt/quest/config"
	"github.com/adetxt/quest/domain"
)

type questUsecase struct {
	cfg       config.Config
	questRepo domain.QuestRepository
}

func NewQuestUsecase(cfg config.Config, questRepo domain.QuestRepository) domain.QuestUsecase {
	return &questUsecase{
		cfg:       cfg,
		questRepo: questRepo,
	}
}

func (u *questUsecase) GetUserQuests(ctx context.Context, params *domain.GetQuestsParam) ([]*domain.QuestUser, error) {
	res, err := u.questRepo.GetUserQuets(ctx, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *questUsecase) AssignQuest(ctx context.Context, questID, userID int64) error {
	if err := u.questRepo.AssignQuest(ctx, questID, userID); err != nil {
		return err
	}

	return nil
}

func (u *questUsecase) RevokeQuest(ctx context.Context, questID, userID int64) error {
	if err := u.questRepo.RevokeQuest(ctx, questID, userID); err != nil {
		return err
	}

	return nil
}

func (u *questUsecase) CompleteObjective(ctx context.Context, objectiveID, userID int64) error {
	if err := u.questRepo.SetObjectiveStatus(ctx, objectiveID, userID, domain.ObjectiveComplete); err != nil {
		return err
	}

	return nil
}
