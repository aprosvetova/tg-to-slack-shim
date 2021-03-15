package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(errorHandler())
	r.GET("/bot:token/:method", apiHandler)
	r.POST("/bot:token/:method", apiHandler)

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
