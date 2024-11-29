package service

import (
	"context"
	gError "goatrobotics/errors"
	"goatrobotics/models"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatRoom struct {
	Clients   sync.Map
	Broadcast chan *models.Message
	Rwmutex   sync.RWMutex
	Join      chan string
	Leave     chan string
	Messages  []*models.Message
}

func NewChatRoomService() *ChatRoom {
	return &ChatRoom{
		Clients:   sync.Map{},
		Broadcast: make(chan *models.Message, 100),
		Join:      make(chan string),
		Leave:     make(chan string),
		Messages:  make([]*models.Message, 0),
	}
}

func (m *ChatRoom) Execute() {
	for {
		select {
		case id := <-m.Join:
			m.Clients.Store(id, struct{}{})
			log.Printf("%s joined the ChatRoom\n", id)
		case id := <-m.Leave:
			m.Clients.LoadAndDelete(id)
			log.Printf("%s Left the ChatRoom\n", id)
		case message := <-m.Broadcast:
			m.Rwmutex.Lock()
			m.Messages = append(m.Messages, &message)
			m.Rwmutex.Unlock()
			log.Printf("Broadcasting Message In the room : %v", message)
		}
	}
}

func (c *ChatRoom) JoinClient(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.Json(http.StatusBadRequest, gin.H{gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); exists {
		log.Println("Client Id is already in Use")
		ctx.Json(http.StatusConflict, gin.H{gError.DUPLICATE_CLIENT_ID})
		return
	}
	c.Join <- id
	ctx.JSON(http.StatusOK, gin.H{models.JoinClientResponse{Id: id, Message: "Joined chat successfully"}})
}

func (c *ChatRoom) LeaveClient(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.Json(http.StatusBadRequest, gin.H{gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.Json(http.StatusConflict, gin.H{gError.CLIENT_ID_NOT_FOUND})
		return
	}

	c.Leave <- id
	ctx.JSON(http.StatusOK, gin.H{models.LeaveClientResponse{Id: id, Message: "Left chat successfully"}})
}

func (c *ChatRoom) SendMessage(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.Json(http.StatusBadRequest, gin.H{gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.Json(http.StatusConflict, gin.H{gError.CLIENT_ID_NOT_FOUND})
		return
	}
	messageText := strings.TrimSpace(ctx.Query("message"))
	if messageText == "" {
		log.Println("Message is Empty")
		ctx.Json(http.StatusBadRequest, gin.H{gError.MESSAGE_IS_EMPTY})
		return

	}

	// message := &models.Message{
	// 	Id:      id,
	// 	Message: messageText,
	// }

	c.Broadcast <- &models.Message{Id: id, Message: messageText}
	ctx.JSON(http.StatusOK, gin.H{models.SendMessageResponse{Id: id, Message: "Message sent successfully"}})
}

func (c *ChatRoom) GetMessages(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.Json(http.StatusBadRequest, gin.H{gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.Json(http.StatusConflict, gin.H{gError.CLIENT_ID_NOT_FOUND})
		return
	}
	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // this may be useful when we wanna do some database kind of stuffs
	defer cancel()

	responseHelper := make(chan []*models.Message, 1)
	go func() {
		c.Rwmutex.RLock()
		messages := make([]*models.Message, len(c.Messages))
		copy(messages, c.Messages)
		c.Rwmutex.RUnlock()

		// Send the fetched messages to the channel
		responseHelper <- messages
	}()

	// Wait for either a response or timeout
	select {
	case messages := <-responseHelper:
		// Send response if messages are retrieved successfully
		response := models.GetMessagesResponse{
			ID:       id,
			Messages: messages,
		}
		if len(messages) == 0 {
			response.Message = "No new messages"
		}

		log.Printf("Info: Sending messages for client ID: %s", id) // Log response
		ctx.JSON(http.StatusOK, response)                          // Send the JSON response back to the client

	case <-requestCtx.Done():
		// Handle timeout
		log.Printf("Error: Timeout retrieving messages for client ID: %s", id) // Log timeout
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Request timed out while retrieving messages",
		})
	}
}
