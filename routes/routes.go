package routes

import (
	"goatrobotics/service"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, chatRoom *service.ChatRoom) {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	router.Use(cors.New(corsConfig))
	router.GET("/join", ValidateIDMiddleware(), chatRoom.JoinClient)
	router.GET("/leave", ValidateIDMiddleware(), chatRoom.LeaveClient)
	router.GET("/send", ValidateIDMiddleware(), chatRoom.SendMessage)
	router.GET("/messages", ValidateIDMiddleware(), chatRoom.GetMessages)
	router.GET("/ping",service.Ping)
}

func ValidateIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
