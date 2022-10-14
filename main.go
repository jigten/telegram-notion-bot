package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var allowedChatIds = map[int]bool{
	-898876211: true,
	1251311320: true,
}

const startCommand string = "/start"

var lenStartCommand int = len(startCommand)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "TELEGRAM_BOT_TOKEN"

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

func parseTelegramRequest(c *gin.Context) (*Update, error) {
	var update Update
	if err := c.BindJSON(&update); err != nil {
		return nil, err
	}

	if update.UpdateId == 0 {
		log.Printf("Invalid update id, got update id = 0")
		return nil, errors.New("Invalid update id of 0 indicates failure to parse incoming update")
	}

	return &update, nil
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	fmt.Printf("BOT TOKEN FROM OS IS %s", os.Getenv("TELEGRAM_BOT_TOKEN"))
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func handler(c *gin.Context) {
	var update, err = parseTelegramRequest(c)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	fmt.Printf("recieved update with contents: %+v\n", update)

	// reject message if not from allowed chats
	if _, ok := allowedChatIds[update.Message.Chat.Id]; !ok {
		sendTextToTelegramChat(update.Message.Chat.Id, "Unauthorized Chat ID")
	}

	// Send response back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, "Hi Tshogyal & Jigten!")
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
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
