package main

import (
	"github.com/TON-Market/tma/server/datatype/market"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"sync"
)

type socket struct {
	wsCh    chan *market.EventDTO
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func newSocket() *socket {
	return &socket{
		wsCh:    market.GetMarket().WsCh,
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.Mutex{},
	}
}

func (h *socket) broadcastEvent(event *market.EventDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		err := websocket.JSON.Send(client, event)
		if err != nil {
			client.Close()
			delete(h.clients, client)
		}
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
		for e := range h.wsCh {
			h.broadcastEvent(e)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
