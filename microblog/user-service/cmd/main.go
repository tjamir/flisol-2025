package main

import (
	"log"

	"github.com/tjamir/flisol-2025/microblog/user-service/internal/handler"
	"github.com/tjamir/flisol-2025/microblog/user-service/internal/model"
	"github.com/tjamir/flisol-2025/microblog/user-service/internal/repository"
)

func main() {
	db := repository.NewPostgresDB()
	if err := db.Raw("SELECT 1").Error; err != nil {
		log.Fatalf("Falha na conexão com o banco: %v", err)
	}
	log.Println("Conexão com o banco estabelecida com sucesso!")
    db.AutoMigrate(&model.User{})
	StartUserService(handler.NewUserHandler(repository.NewPostgresUserRepository(db)))
}