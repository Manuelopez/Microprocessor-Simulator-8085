package main

import (
	"encoding/json"
	"fmt"
	"log"
	"micp-sim/microprocessor"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var clients map[*Client]bool = make(map[*Client]bool)

type WsData struct {
	Type string `json:"type"`
  Data json.RawMessage `json:"data"`
}

type StartData struct {
  Instructions []string
}



type Client struct {
	Conn *websocket.Conn
	microprocessor.MicroProcessor
}

func (c *Client) reader() {
	defer func() {
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1024)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg WsData

		json.Unmarshal(message, &msg)
		switch msg.Type {
		case "connecting":
			fmt.Println("connected")
		case "start":
      var data StartData
      json.Unmarshal(msg.Data, &data)
      go c.MicroProcessor.Start(data.Instructions, messageType)
      
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Client Connected")

	client := Client{Conn: ws, MicroProcessor: microprocessor.New(5, ws)}

	clients[&client] = true
	go client.reader()

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
	//  m := microprocessor.New(3, nil)
	//  m.Start([]string{"BEGIN", "MOV AL 5", "END"})
}
