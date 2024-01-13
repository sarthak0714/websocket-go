package main

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type Server struct {
	Clients map[*Client]bool
}

func NewServer() *Server {
	return &Server{
		make(map[*Client]bool),
	}
}

func (s *Server) HandleWs(ws *websocket.Conn) {
	defer ws.Close()
	client := &Client{
		Conn: ws,
	}
	s.Clients[client] = true
	s.HandleMessages(client)
}

func (s *Server) HandleMessages(client *Client) {
	for {
		var msg string
		err := websocket.Message.Receive(client.Conn, &msg)
		if err != nil {
			if err == io.EOF {
				delete(s.Clients, client)
			}
			break
		}
		for client := range s.Clients {
			err := websocket.Message.Send(client.Conn, msg)
			if err != nil {
				delete(s.Clients, client)
			}
		}
	}
}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.HandleWs))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err.Error())
	}
}
