package questmysql

import (
	"context"
	"sync"
	"time"

	"github.com/adetxt/quest/config"
	"github.com/adetxt/quest/domain"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type repository struct {
	cfg config.Config
	db  *gorm.DB
}

func New(cfg config.Config, db *gorm.DB) domain.QuestRepository {
	return &repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *repository) GetUserQuets(ctx context.Context, params *domain.GetQuestsParam) ([]*domain.QuestUser, error) {
	questUsers := []QuestUser{}

	query := r.db.Debug().Where("quest_users.user_id = ?", params.UserID)

	if params.Status == domain.ObjectiveComplete.String() {
		query.Joins("JOIN objectives ON objectives.quest_id = quest_users.quest_id").
			Joins("JOIN objective_users ON objective_users.objective_id = objectives.id").
			Group("quest_users.quest_id").Having("SUM(IF(objective_users.status = 'IN_PROGRESS', 1, 0)) < 1")
	}

	if err := query.Find(&questUsers).Error; err != nil {
		return nil, err
	}

	questIDs := []int64{}
	for i := 0; i < len(questUsers); i++ {
		questIDs = append(questIDs, questUsers[i].QuestID)
	}

	eg, _ := errgroup.WithContext(ctx)

	// group by quest id
	objectivesChan := make(chan map[int64][]domain.Objective)
	eg.Go(func() error {
		defer close(objectivesChan)
		objectives := []Objective{}
		if err := r.db.Model(&Objective{}).Where("quest_id IN ?", questIDs).Find(&objectives).Error; err != nil {
			return err
		}

		res := map[int64][]domain.Objective{}
		for _, v := range objectives {
			res[v.QuestID] = append(res[v.QuestID], *v.ToEntity())
		}

		objectivesChan <- res
		return nil
	})

	// key by quest id
	questsChan := make(chan map[int64]Quest)
	eg.Go(func() error {
		defer close(questsChan)
		quests := []Quest{}
		if err := r.db.Model(&Quest{}).Where("id IN ?", questIDs).Find(&quests).Error; err != nil {
			return err
		}

		res := map[int64]Quest{}
		for _, v := range quests {
			res[v.ID] = v
		}

		questsChan <- res
		return nil
	})

	// map by objective id
	objectiveUsersChan := make(chan map[int64]ObjectiveUser)
	eg.Go(func() error {
		defer close(objectiveUsersChan)
		objectiveUsers := []ObjectiveUser{}
		if err := r.db.Model(&ObjectiveUser{}).Where("user_id = ?", params.UserID).Find(&objectiveUsers).Error; err != nil {
			return err
		}

		res := map[int64]ObjectiveUser{}
		for _, v := range objectiveUsers {
			res[v.ObjectiveID] = v
		}

		objectiveUsersChan <- res
		return nil
	})

	objectives := <-objectivesChan
	quests := <-questsChan
	objectiveUsers := <-objectiveUsersChan

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	res := make([]*domain.QuestUser, len(questUsers))

	var wg sync.WaitGroup

	for i := 0; i < len(questUsers); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			q := quests[questUsers[i].QuestID]

			res[i] = &domain.QuestUser{
				ID:          q.ID,
				Title:       q.Title,
				Description: q.Description,
			}

			totalComplete := 0
			for _, v := range objectives[q.ID] {
				status := domain.ObjectiveInprogress

				if obj, ok := objectiveUsers[v.ID]; ok && obj.Status == domain.ObjectiveComplete.String() {
					totalComplete++
					status = domain.ObjectiveStatus(obj.Status)
				}

				res[i].Objectives = append(res[i].Objectives, domain.ObjectiveUser{
					ID:          v.ID,
					QuestID:     q.ID,
					Title:       v.Title,
					Description: v.Description,
					Status:      status,
				})

				res[i].Percentage = (float32(totalComplete) / float32(len(objectives[q.ID]))) * 100
			}
		}(i)
	}

	wg.Wait()

	return res, nil
}

func (r *repository) AssignQuest(ctx context.Context, questID, userID int64) error {
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		quest := Quest{}
		if err := r.db.Where("id = ?", questID).First(&quest).Error; err != nil {
			return err
		}

		return nil
	})

	objectiveUsersChan := make(chan []ObjectiveUser)
	eg.Go(func() error {
		defer close(objectiveUsersChan)

		objectives := []Objective{}
		if err := r.db.Where("quest_id = ?", questID).Find(&objectives).Error; err != nil {
			return err
		}

		objectiveUsers := []ObjectiveUser{}
		for _, v := range objectives {
			objectiveUsers = append(objectiveUsers, ObjectiveUser{
				UserID:      userID,
				ObjectiveID: v.ID,
				Status:      domain.ObjectiveInprogress.String(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
		}

		objectiveUsersChan <- objectiveUsers
		return nil
	})

	objectiveUsers := <-objectiveUsersChan

	if err := eg.Wait(); err != nil {
		return err
	}

	return r.db.Debug().Transaction(func(tx *gorm.DB) error {
		q := QuestUser{
			UserID:    userID,
			QuestID:   questID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&q).Error; err != nil {
			return err
		}

		if err := tx.Create(&objectiveUsers).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) RevokeQuest(ctx context.Context, questID, userID int64) error {
	objectives := []Objective{}
	if err := r.db.Where("quest_id = ?", questID).Find(&objectives).Error; err != nil {
		return err
	}

	objectiveIDs := []int64{}
	for _, v := range objectives {
		objectiveIDs = append(objectiveIDs, v.ID)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND quest_id = ?", userID, questID).
			Delete(&QuestUser{}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ? AND objective_id IN ?", userID, objectiveIDs).
			Delete(&ObjectiveUser{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) SetObjectiveStatus(ctx context.Context, objectiveID, userID int64, status domain.ObjectiveStatus) error {
	if err := r.db.Table("objective_users").Where("user_id = ? AND objective_id = ?", userID, objectiveID).
		Update("status", domain.ObjectiveComplete).Error; err != nil {
		return err
	}

	return nil
}
