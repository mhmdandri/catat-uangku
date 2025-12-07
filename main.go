package main

import (
	"catatan-keuangan/config"
	"catatan-keuangan/database"
	"catatan-keuangan/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	routes.InitRoutes(r)
	log.Println("Server berjalan di", config.Cfg.AppPort)
	if err := r.Run(config.Cfg.AppPort); err != nil {
		log.Fatal("Gagal menjalankan server", err)
	}
}
