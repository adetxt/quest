package usermysql

import "github.com/adetxt/quest/domain"

type User struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Avatar   string `gorm:"column:avatar"`
}

func (i *User) ToEntity() *domain.User {
	return &domain.User{
		ID:       i.ID,
		Username: i.Username,
		Email:    i.Email,
		Avatar:   i.Avatar,
	}
}

func (i *User) FromEntity(d *domain.User) {
	i.ID = d.ID
	i.Username = d.Username
	i.Email = d.Email
	i.Avatar = d.Avatar
}
