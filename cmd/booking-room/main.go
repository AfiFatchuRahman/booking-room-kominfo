package main

import (
	"log"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/config"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/handler"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/repository"
	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	db, err := config.OpenDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("database init error: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	uc := usecase.NewUsecase(repo, cfg.JWTSecret)
	h := handler.NewHTTPHandler(uc)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	h.RegisterRoutes(r, cfg.JWTSecret)

	log.Printf("server berjalan pada port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server gagal dijalankan: %v", err)
	}
}
