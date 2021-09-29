package main

import (
	"bytes"
	"chat-api/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	ChannelId   string    `json:"channelId"`
	Message     string    `json:"message"`
	MessageFrom string    `json:"messageFrom"`
	TimeStamp   time.Time `json:"timeStamp"`
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			for client := range h.clients {
				if err := client.WriteJSON(message); !errors.Is(err, nil) {
					log.Printf("error occurred: %v", err)
				}
			}
		}
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Create a hub
	hub := NewHub()

	// Start a go routine
	go hub.run()
	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer func() {
			delete(hub.clients, ws)
			ws.Close()
			log.Printf("Closed!")
		}()

		// Add client
		hub.clients[ws] = true

		log.Println("Connected!")

		// Listen on connection
		read(hub, ws)
		return nil
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func read(hub *Hub, client *websocket.Conn) {

	for {
		var message Message
		httpClient := &http.Client{}
		err := client.ReadJSON(&message)

		if !errors.Is(err, nil) {
			log.Printf("error occurred: %v", err)
			delete(hub.clients, client)
			break
		}

		message.TimeStamp = time.Now()
		requestBody, _ := json.Marshal(message)
		req, err := http.NewRequest(http.MethodPost, os.Getenv("SERVER_MESSAGE_URL"), bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		result := usecase.CreateMessageOutput{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
		// Send a message to hub
		hub.broadcast <- message
	}
}
