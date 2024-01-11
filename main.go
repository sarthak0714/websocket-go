package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conn map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWs(ws *websocket.Conn) {
	fmt.Println("new conn from ", ws.RemoteAddr())
	s.conn[ws] = true
	s.rLoop(ws)
}

func (s *Server) rLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed from ", ws.LocalAddr())
				break
			}
			fmt.Println("Err occured with message:", err)
			continue
		}
		msg := buff[:n]
		go s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conn {
		_, err := ws.Write(b)
		if err != nil {
			fmt.Println("err while writing:", err)
		}
	}
}

func main() {
	s := NewServer()
	http.Handle("/chat", websocket.Handler(s.HandleWs))
	http.ListenAndServe(":3000", nil)
}
