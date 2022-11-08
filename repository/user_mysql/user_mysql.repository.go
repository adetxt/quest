package usermysql

import (
	"context"

	"github.com/adetxt/quest/config"
	"github.com/adetxt/quest/domain"
	"gorm.io/gorm"
)

type repository struct {
	cfg config.Config
	db  *gorm.DB
}

func New(cfg config.Config, db *gorm.DB) domain.UserRepository {
	return &repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *repository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users := []User{}

	if err := r.db.Model(&User{}).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))

	for i := 0; i < len(users); i++ {
		result[i] = users[i].ToEntity()
	}

	return result, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := User{}

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return user.ToEntity(), nil
}
