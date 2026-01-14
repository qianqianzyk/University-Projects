package main

import (
	"SGMS/app/midwares"
	"SGMS/config/database"
	"SGMS/config/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	database.Init()
	r := gin.Default()
	r.Use(midwares.Corss())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	router.Init(r)

	log.Println("Starting server on :8888")
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
