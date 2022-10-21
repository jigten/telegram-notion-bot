package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	telegram_service "github.com/jigten/telegram-notion-bot/services/telegram_service"
	command_handler "github.com/jigten/telegram-notion-bot/util/command_handler"
	greeting "github.com/jigten/telegram-notion-bot/util/greeting"
	"github.com/joho/godotenv"
)

var allowedChatIds = map[int]bool{
	-898876211: true,
	1251311320: true,
}

func handler(c *gin.Context) {
	var update, err = telegram_service.ParseTelegramRequest(c)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	fmt.Printf("recieved update with contents: %+v\n", update)

	// reject message if not from allowed chats
	if _, ok := allowedChatIds[update.Message.Chat.Id]; !ok {
		telegram_service.SendTextToTelegramChat(update.Message.Chat.Id, "Unauthorized Chat ID")
	}

	command, _, parseErr := telegram_service.ParseEventCommand(update.Message.Text)
	if parseErr != nil {
		telegram_service.SendTextToTelegramChat(update.Message.Chat.Id, parseErr.Error())
	}

	returnMsg, msgErr := command_handler.HandleCommand(command)
	if msgErr != nil {
		telegram_service.SendTextToTelegramChat(update.Message.Chat.Id, msgErr.Error())
	}

	// Send response back to Telegram
	telegramResponseBody, errTelegram := telegram_service.SendTextToTelegramChat(update.Message.Chat.Id, returnMsg)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("successfully distributed to chat id %d", update.Message.Chat.Id)
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/handler", handler)

	return r
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	// r.Run(":8080")
	greeting.SetGreeting("Hello this is a new greeting")
	greeting := greeting.ReadGreetingFile()
	fmt.Printf("%s", greeting)
}
