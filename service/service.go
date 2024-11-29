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
			m.Messages = append(m.Messages, message)
			m.Rwmutex.Unlock()
			log.Printf("Broadcasting Message In the room : %v", message)
		}
	}
}

func (m *ChatRoom) JoinClient(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error":gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); exists {
		log.Println("Client Id is already in Use")
		ctx.JSON(http.StatusConflict, gin.H{"error":gError.DUPLICATE_CLIENT_ID})
		return
	}
	m.Join <- id
	ctx.JSON(http.StatusOK, models.JoinClientResponse{Id: id, Message: "Joined chat successfully"})
}

func (m *ChatRoom) LeaveClient(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error":gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.JSON(http.StatusConflict, gin.H{"error":gError.CLIENT_ID_NOT_FOUND})
		return
	}

	m.Leave <- id
	ctx.JSON(http.StatusOK, models.LeaveClientResponse{Id: id, Message: "Left chat successfully"})
}

func (m *ChatRoom) SendMessage(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error":gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.JSON(http.StatusConflict, gin.H{"error":gError.CLIENT_ID_NOT_FOUND})
		return
	}
	messageText := strings.TrimSpace(ctx.Query("message"))
	if messageText == "" {
		log.Println("Message is Empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error":gError.MESSAGE_IS_EMPTY})
		return

	}

	m.Broadcast <- &models.Message{Id: id, Message: messageText}
	ctx.JSON(http.StatusOK, models.SendMessageResponse{Id: id, Message: "Message sent successfully"})
}

func (m *ChatRoom) GetMessages(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Query("id"))
	if id == "" {
		log.Println("Client Id is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error":gError.CLIENT_ID_REQUIRED})
		return
	}
	if _, exists := m.Clients.Load(id); !exists {
		log.Println("Client Id is Not Present In the Room")
		ctx.JSON(http.StatusConflict, gin.H{"error":gError.CLIENT_ID_NOT_FOUND})
		return
	}
	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second) 
	defer cancel()

	responseHelper := make(chan []*models.Message, 1)
	go func() {
		m.Rwmutex.RLock()
		messages := make([]*models.Message, len(m.Messages))
		copy(messages, m.Messages)
		m.Rwmutex.RUnlock()
		responseHelper <- messages
	}()
	select {
	case messages := <-responseHelper:
		response := models.GetMessagesResponse{
			Id:       id,
			Messages: messages,
		}
		if len(messages) == 0 {
			response.MessageIndicator = "No new messages"
		}

		log.Printf("Info: Sending messages for client ID: %s\n", id) 
		ctx.JSON(http.StatusOK, response)                         

	case <-requestCtx.Done():
		
		log.Printf("Error: Timeout retrieving messages for client ID: %s", id) 
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Request timed out while retrieving messages",
		})
	}
}


func Ping(ctx *gin.Context) {
	log.Println("[INFO] Received ping request")
	response := &models.PingResponse{Message: "Pinged Successfully"}
	ctx.JSON(200, response)
}