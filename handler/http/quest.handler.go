package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adetxt/quest/domain"
	"github.com/adetxt/quest/utils/edison"
)

type (
	assignQuestRequest struct {
		QuestID int64 `json:"quest_id"`
		UserID  int64 `json:"user_id"`
	}

	completeObjectiveRequest struct {
		ObjectiveID int64 `json:"objective_id"`
		UserID      int64 `json:"user_id"`
	}
)

type questHandler struct {
	edison  *edison.Edison
	questUc domain.QuestUsecase
}

func NewQuestHandler(ed *edison.Edison, questUc domain.QuestUsecase) {
	h := &questHandler{ed, questUc}

	ed.RestRouter("GET", "/api/v1/quests/user/:userID", h.GetUserQuests)
	ed.RestRouter("POST", "/api/v1/quests/assign", h.AssignQuest)
	ed.RestRouter("POST", "/api/v1/quests/revoke", h.RevokeQuest)
	ed.RestRouter("POST", "/api/v1/quests/objective/complete", h.CompleteObjective)
}

func (h *questHandler) GetUserQuests(ctx context.Context, c edison.RestContext) error {
	userID, err := strconv.Atoi(c.EchoContext.Param("userID"))
	if err != nil {
		return c.Error(err, http.StatusBadRequest)
	}

	status := ""
	if c.EchoContext.QueryParam("only_complete") == "true" {
		status = domain.ObjectiveComplete.String()
	}

	fmt.Printf("status %s", status)

	users, err := h.questUc.GetUserQuests(ctx, &domain.GetQuestsParam{
		UserID: int64(userID),
		Status: status,
	})
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	return c.Success(users)
}

func (h *questHandler) AssignQuest(ctx context.Context, c edison.RestContext) error {
	req := &assignQuestRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if req.QuestID < 1 || req.UserID < 1 {
		return c.ErrorWithCustomMessage(http.StatusBadRequest, "quest_id and user_id is required")
	}

	err := h.questUc.AssignQuest(ctx, req.QuestID, req.UserID)
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	return c.Success("ok")
}

func (h *questHandler) RevokeQuest(ctx context.Context, c edison.RestContext) error {
	req := &assignQuestRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if req.QuestID < 1 || req.UserID < 1 {
		return c.ErrorWithCustomMessage(http.StatusBadRequest, "quest_id and user_id is required")
	}

	err := h.questUc.RevokeQuest(ctx, req.QuestID, req.UserID)
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	return c.Success("ok")
}

func (h *questHandler) CompleteObjective(ctx context.Context, c edison.RestContext) error {
	req := &completeObjectiveRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if req.ObjectiveID < 1 || req.UserID < 1 {
		return c.ErrorWithCustomMessage(http.StatusBadRequest, "objective_id and user_id is required")
	}

	err := h.questUc.CompleteObjective(ctx, req.ObjectiveID, req.UserID)
	if err != nil {
		return c.Error(err, http.StatusInternalServerError)
	}

	return c.Success("ok")
}
