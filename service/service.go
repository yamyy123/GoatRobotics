package service

import (
	"goatrobotics/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatRoom struct {
	Clients   map[string]chan *models.Message
	Broadcast chan *models.Message
	Rwmutex   sync.RWMutex
	Join      chan string
	Leave     chan string
	Messages  []*models.Message
}

func NewChatRoomService() *ChatRoom {
	return &ChatRoom{
		Clients:   make(map[string]chan *models.Message, 100),
		Broadcast: make(chan *models.Message),
		Join:      make(chan string),
		Leave:     make(chan string),
		Messages:  make([]*models.Message, 0),
	}
}

func (m *ChatRoom) Execute() {
	for {
		select {
		case id := <-m.Join:
			m.Rwmutex.Lock()
			if _, exists := m.Clients[id]; !exists {
				m.Clients[id] = make(chan *models.Message, 100)
			}
			m.Rwmutex.Unlock()
		case id := <-m.Leave:
			m.Rwmutex.Lock()
			if clientChan, exists := m.Clients[id]; exists {
				close(clientChan)
				delete(m.Clients, id)
			}
			m.Rwmutex.Unlock()
		case message := <-m.Broadcast:
			m.Rwmutex.RLock()
			for id, clientChan := range m.Clients {
				go func(id string, clientChan chan *models.Message) {
					select {
					case clientChan <- message: // Try sending the message
						log.Printf("Message sent to client: %s", id)
					default:
						log.Printf("Client %s is not consuming messages", id)
					}
				}(id, clientChan)
			}
			m.Rwmutex.RUnlock()
		}
	}
}



func (c *ChatRoom) JoinClient(ctx *gin.Context) {
	id := ctx.Query("id")
	c.Join <- id
	ctx.JSON(http.StatusOK, gin.H{"message": "Joined chat successfully", "id": id})
}

// LeaveClient handles the leave chat request
func (c *ChatRoom) LeaveClient(ctx *gin.Context) {
	id := ctx.Query("id")
	c.Leave <- id
	ctx.JSON(http.StatusOK, gin.H{"message": "Left chat successfully", "id": id})
}

// SendMessage handles sending a message to the chat room
func (c *ChatRoom) SendMessage(ctx *gin.Context) {
	id := ctx.Query("id")
	messageText := ctx.Query("message")

	message := &models.Message{
		Id:      id,
		Message: messageText,
	}

	c.Broadcast <- message
	ctx.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "id": id})
}

// GetMessages handles retrieving messages for a client
func (c *ChatRoom) GetMessages(ctx *gin.Context) {
	id := ctx.Query("id")

	c.Rwmutex.RLock()
	clientChan, exists := c.Clients[id]
	c.Rwmutex.RUnlock()
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	select {
	case message := <-clientChan:
		ctx.JSON(http.StatusOK, message)
	case <-time.After(3 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "No new messages"})
	}
}
