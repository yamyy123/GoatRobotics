package routes

import (
	"bytes"
	"fmt"
	"goatrobotics/service"
	"io"
	"log"
	"os"
	"time"

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
    router.Static("/static", "./UI")
	router.Use(cors.New(corsConfig))
	router.GET("/join", AuditMiddleware(), chatRoom.JoinClient)
	router.GET("/leave", AuditMiddleware(), chatRoom.LeaveClient)
	router.GET("/send", AuditMiddleware(), chatRoom.SendMessage)
	router.GET("/messages", AuditMiddleware(), chatRoom.GetMessages)
	router.GET("/ping", AuditMiddleware(), service.Ping)

}

func AuditMiddleware() gin.HandlerFunc {
	logFile, err := os.OpenFile("audits/audit_logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	return func(ctx *gin.Context) {
		startTime := time.Now()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		headers := ctx.Request.Header
		userAgent := ctx.Request.UserAgent()
		path := ctx.FullPath()
		queryParams := ctx.Request.URL.Query()
		hostname, _ := os.Hostname()
		var requestBody bytes.Buffer
		ctx.Request.Body = io.NopCloser(io.TeeReader(ctx.Request.Body, &requestBody))
		requestData := requestBody.String()

		var responseBody bytes.Buffer
		writer := &bodyWriter{body: &responseBody, ResponseWriter: ctx.Writer}
		ctx.Writer = writer
		ctx.Next()
		responseCode := ctx.Writer.Status()
		responseData := responseBody.String()
		responseSize := writer.body.Len()
		duration := time.Since(startTime)
		logEntry := fmt.Sprintf(
			"Time: %s | Hostname: %s | ClientIP: %s | Method: %s | Path: %s | QueryParams: %v | Headers: %v | UserAgent: %s | RequestBody: %s | ResponseCode: %d | ResponseSize: %d | Response: %s | Duration: %v\n",
			time.Now().Format(time.RFC3339), hostname, clientIP, method, path, queryParams, headers, userAgent, requestData, responseCode, responseSize, responseData, duration,
		)

		logger.Println(logEntry)
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
