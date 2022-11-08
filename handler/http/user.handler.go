package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adetxt/quest/domain"
	"github.com/adetxt/quest/utils/edison"
)

// type (
// 	createDatabaseRequest struct {
// 		Service string `json:"service"`
// 		Host    string `json:"host"`
// 		Port    string `json:"port"`
// 		DBName  string `json:"db_name"`
// 	}
// )

type userHandler struct {
	edison *edison.Edison
	userUc domain.UserUsecase
}

func NewUserHandler(ed *edison.Edison, userUc domain.UserUsecase) {
	h := &userHandler{ed, userUc}

	ed.RestRouter("GET", "/api/v1/users", h.GetUsers)
	ed.RestRouter("GET", "/api/v1/users/:username", h.GetUser)
}

func (h *userHandler) GetUsers(ctx context.Context, c edison.RestContext) error {
	users, err := h.userUc.GetUsers(ctx)
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	return c.Success(users)
}

func (h *userHandler) GetUser(ctx context.Context, c edison.RestContext) error {
	user, err := h.userUc.GetUserByIdentifier(ctx, "username", c.EchoContext.Param("username"))
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	if user == nil {
		return c.Error(fmt.Errorf("not found"), http.StatusNotFound)
	}

	return c.Success(user)
}
