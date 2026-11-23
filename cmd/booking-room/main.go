package main

import (
	"log"

	"github.com/AfiFatchuRahman/booking-room-kominfo/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	db, err := config.OpenDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("database init error: %v", err)
	}
	defer db.Close()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	log.Printf("server berjalan pada port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server gagal dijalankan: %v", err)
	}
}
