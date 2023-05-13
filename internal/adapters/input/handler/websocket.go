// Package handler implements the handlers
package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/posilva/simplechat/internal/core/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebSockerHandler manages the handling of HTTP requests
type WebSockerHandler struct {
	chat services.ChatService
}

// NewWebSockerHandler creates a new HTTP handler
func NewWebSockerHandler(chat services.ChatService) *WebSockerHandler {
	return &WebSockerHandler{
		chat: chat,
	}
}

// Handle handles the request to send message
func (h *WebSockerHandler) Handle(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// TODO: we may want to not return the error as it can contain info we do not want to share
		// 			instead we may want to write to logs or send some metric
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to upgrade: %v", err))
	}
	defer func() {
		_ = ws.Close()
	}()

	// TODO: this will require better session management
	id := ctx.Query("id")
	room := ctx.Query("room")

	rc := newclientReceiver()
	ep := newclientEndpoint(id, room, rc)

	err = h.chat.Register(ep)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to register client in the chat service: %v", err))
		return
	}
	go func() {
		defer func() {
			_ = h.chat.UnRegister(ep)
			_ = ws.Close()
		}()

		for {
			_, b, err := ws.ReadMessage()
			if err != nil {
				_, _ = fmt.Printf("failed to read message from socket: %v", err)
				break
			}

			m, err := decode(b)
			if err != nil {
				_, _ = fmt.Printf("failed to decode message: %v", err)
				break
			}
			err = h.chat.Send(m)
			if err != nil {
				_, _ = fmt.Printf("failed to send message to chat service: %v", err)
				break
			}
		}
	}()

	defer func() {
		_ = h.chat.UnRegister(ep)
	}()

	for mm := range rc.Channel() {
		m, err := encode(mm.Message)
		if err != nil {
			_, _ = fmt.Printf("failed to encode message: %v", err)
			break
		}
		err = ws.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			_, _ = fmt.Printf("failed to write to socket: %v", err)
			break
		}
	}
}
