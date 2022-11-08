package main

import (
	"github.com/adetxt/quest/config"
	httpHdl "github.com/adetxt/quest/handler/http"
	questmysql "github.com/adetxt/quest/repository/quest_mysql"
	usermysql "github.com/adetxt/quest/repository/user_mysql"
	"github.com/adetxt/quest/usecase"
	"github.com/adetxt/quest/utils/edison"
	"github.com/adetxt/quest/utils/mysql"
)

func main() {
	cfg := config.New()

	// external agencies
	db := mysql.Init(&mysql.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		DBName:   cfg.DBName,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
	})

	// !! DEVELOPMENT ONLY
	db.AutoMigrate(usermysql.User{}, questmysql.Quest{}, questmysql.Objective{}, questmysql.QuestUser{}, questmysql.ObjectiveUser{})

	// Repositories
	userRepo := usermysql.New(cfg, db)
	questRepo := questmysql.New(cfg, db)

	// Usecases
	userUc := usecase.NewUserUsecase(cfg, userRepo, questRepo)
	questUc := usecase.NewQuestUsecase(cfg, questRepo)

	ed := edison.NewEdison()

	// Handler
	httpHdl.NewUserHandler(ed, userUc)
	httpHdl.NewQuestHandler(ed, questUc)

	ed.StartRestServer(cfg.REST_PORT)
}
