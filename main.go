package main

import (
	"fmt"

	"goatrobotics/constants"
	"goatrobotics/routes"
	"goatrobotics/service"
	"goatrobotics/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	if err := utils.LoadConfig(); err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}

	port := viper.GetString("port")
	if port == "" {
		fmt.Println("Port not specified in the configuration, defaulting to port in constants")
		port = constants.Port
	}

	router := gin.Default()
	chatRoom := service.NewChatRoomService()
	go chatRoom.Execute()
	routes.RegisterRoutes(router, chatRoom)

	fmt.Printf("Chat server running on http://localhost:%s\n", port)
	router.Run(":" + port)
}
