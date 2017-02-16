package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {
	// socket is the web socket for this client.
	name   string
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan []byte
	// room is the room this client is chatting in.
	room *room
}

// Read continuously from the socket of the connection.
func (c *client) read() {

	// Ensure the resource is closed by deferring a call to `Close()`
	defer c.socket.Close()

	//Continuously read messages for the life of the function. If an error is encountered, break from the function
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {

			return
		}
		c.room.forward <- msg
	}
}

// Write anything from the send channel be written to the websocket for the client.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {

			return
		}
	}
}
