package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {

	// Connect to the websocket server
	url := "ws://localhost:8080/list"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// Send a JSON-RPC request to add an item to the list
	request := `
	{
		"jsonrpc": "2.0",
		"method": "addListItem",
		"params": {
			"listItem": "apple"
		},
		"id": "1"
	}
	`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		log.Fatal(err)
	}

	// Read the response from the server
	_, response, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	request = `
	{
		"jsonrpc": "2.0",
		"method": "getList",
		"id": "2"
	}
	`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		log.Fatal(err)
	}

	// Read the response from the server
	_, response, err = conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Get List Response: %s\n", string(response))

	// Send a JSON-RPC request to get the status
	request = `
	{
		"jsonrpc": "2.0",
		"method": "status",
		"id": "3"
	}
	`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		log.Fatal(err)
	}

	// Read the response from the server
	_, response, err = conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status Response: %s\n", string(response))

	//Before we can receive notifications, we have to use the login method
	request = `
	{
		"jsonrpc": "2.0",
		"method": "login",
		"params": {
			"username": "test"
		},
		"id": "4"
	}
	`

	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		log.Fatal(err)
	}
	// Read the response from the server
	_, response, err = conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Login Response: %s\n", string(response))

	// Continuously read notifications from the server
	for {
		_, notification, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received Notification: %s\n", string(notification))
	}

}
