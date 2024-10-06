package main

import (
	"github.com/TON-Market/tma/server/datatype/event"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"sync"
)

type socket struct {
	clients   map[*websocket.Conn]bool // Активные WebSocket соединения
	mu        sync.Mutex               // Для безопасного параллельного доступа к клиентам
	eventChan chan *event.EventDTO     // Канал для событий
}

func newSocket() *socket {
	return &socket{
		clients:   make(map[*websocket.Conn]bool),
		eventChan: event.Keeper().Ch, // Предполагаем, что у тебя есть структура Event
	}
}

// Метод для отправки события всем клиентам
func (h *socket) broadcastEvent(event *event.EventDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Проходим по всем клиентам и отправляем событие
	for client := range h.clients {
		err := websocket.JSON.Send(client, event)
		if err != nil {
			// Если произошла ошибка при отправке, закрываем соединение и удаляем клиента
			client.Close()
			delete(h.clients, client)
		}
	}
}

// WebSocket обработчик
func (h *socket) updateEvent(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		// Добавляем нового клиента
		h.mu.Lock()
		h.clients[ws] = true
		h.mu.Unlock()

		defer func() {
			// Удаляем клиента при закрытии соединения
			h.mu.Lock()
			delete(h.clients, ws)
			h.mu.Unlock()
			ws.Close()
		}()

		// Чтение данных из канала событий и рассылка клиентам
		for event := range h.eventChan {
			h.broadcastEvent(event)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
