package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var globalRoom *room

var roomStore map[string]*room

type room struct {
	//the name of the channel
	name string
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *client

	orders chan *Order
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
}

func newRoom(name string) *room {
	if name == "" {
		return &room{
			name:    "ignore",
			forward: make(chan []byte),
			join:    make(chan *client),
			orders:  make(chan *Order),
			leave:   make(chan *client),
			clients: make(map[*client]bool),
		}
	}
	return &room{
		name:    name,
		forward: make(chan []byte),
		join:    make(chan *client),
		orders:  make(chan *Order),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining

			r.clients[client] = true

		case order := <-r.orders:
			for client := range r.clients {
				client.send <- []byte(order.UUID)
			}
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func handleFromRoom(r *room) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		client := &client{
			socket: socket,
			send:   make(chan []byte, messageBufferSize),
			room:   r,
		}
		r.join <- client

		defer func() { r.leave <- client }()
		go client.write()
		client.read()
		log.Println("Done Reading")

	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize, CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		Error.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	ip := req.RemoteAddr
	r.join <- client
	r.forward <- []byte(ip + " joined the room")
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
