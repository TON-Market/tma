package main

import (
	"log"
	"sync"

	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type socket struct {
	wsCh    chan *market.EventDTO
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func newSocket() *socket {
	return &socket{
		wsCh:    make(chan *market.EventDTO, 100),
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *socket) broadcastEvent(event *market.EventDTO) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		go func(c *websocket.Conn) {
			err := websocket.JSON.Send(c, event)
			if err != nil {
				log.Printf("Ошибка отправки клиенту: %v", err)
				h.mu.Lock()
				c.Close()
				delete(h.clients, c)
				h.mu.Unlock()
			}
		}(client)
	}
}

func (h *socket) updateEvent(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		h.mu.Lock()
		h.clients[ws] = true
		h.mu.Unlock()

		defer func() {
			h.mu.Lock()
			delete(h.clients, ws)
			h.mu.Unlock()
			ws.Close()
		}()

		for event := range h.wsCh {
			h.broadcastEvent(event)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
