package main

import (
	"catatan-keuangan/config"
	"catatan-keuangan/database"
	"catatan-keuangan/middleware"
	"catatan-keuangan/routes"
	"log"

	_ "catatan-keuangan/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title CatatUangku API
// @version 1.0
// @description Swagger documentation untuk layanan CatatUangku. Import file JSON-nya ke Postman untuk koleksi otomatis.
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadConfig()
	database.ConnectDB()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Static("/uploads", "./uploads")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.InitRoutes(r)
	log.Println("Server berjalan di", config.Cfg.AppPort)
	if err := r.Run(config.Cfg.AppPort); err != nil {
		log.Fatal("Gagal menjalankan server", err)
	}
}
