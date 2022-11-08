package usecase

import (
	"context"
	"errors"

	"github.com/adetxt/quest/config"
	"github.com/adetxt/quest/domain"
	"gorm.io/gorm"
)

type userUsecase struct {
	cfg       config.Config
	userRepo  domain.UserRepository
	questRepo domain.QuestRepository
}

func NewUserUsecase(cfg config.Config, userRepo domain.UserRepository, questRepo domain.QuestRepository) domain.UserUsecase {
	return &userUsecase{
		cfg:       cfg,
		userRepo:  userRepo,
		questRepo: questRepo,
	}
}

func (u *userUsecase) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := u.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, nil
	}

	return users, nil
}

func (u *userUsecase) GetUserByIdentifier(ctx context.Context, identifier string, value string) (user *domain.User, err error) {
	switch identifier {
	case "username":
		user, err = u.userRepo.GetUserByUsername(ctx, value)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, nil
}
