package router

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type WebSocket struct {
	mu        sync.Mutex
	conn      *websocket.Conn
	listeners []webSocketListener

	close  chan struct{}
	closed chan struct{}

	incoming chan []byte
	outgoing chan []byte
}

type WebSocketOptions struct {
	PingTime        time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

type webSocketListener interface {
	onWebSocketOpened(ws *WebSocket)
	onWebSocketClosed(ws *WebSocket)
}

func NewWebSocket(w http.ResponseWriter, r *http.Request, opts *WebSocketOptions) (*WebSocket, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ws := &WebSocket{
		conn:     conn,
		close:    make(chan struct{}, 1),
		closed:   make(chan struct{}),
		incoming: make(chan []byte, 8),
		outgoing: make(chan []byte, 8),
	}

	go ws.mainLoop()
	go ws.mainLoop()
	go ws.readLoop()
	return ws, nil
}

func (ws *WebSocket) Closed() <-chan struct{} { return ws.closed }
func (ws *WebSocket) Incoming() <-chan []byte { return ws.incoming }
func (ws *WebSocket) Outgoing() chan<- []byte { return ws.outgoing }

func (ws *WebSocket) Close() <-chan struct{} {
	select {
	case ws.close <- struct{}{}:
	default:
	}
	return ws.closed
}

func (ws *WebSocket) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			ws.Close()
			// TODO: Log a panic.
		}
	}()

	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			// TODO: Log a read error.
			return
		}
		// TODO: Log an incoming message.

		select {
		case ws.incoming <- msg:
		case <-ws.closed:
		}
	}
}

func (ws *WebSocket) mainLoop() {
	defer func() {
		if err := recover(); err != nil {
			// TODO: Log a panic.
		}
	}()

	defer ws.onClosed()
	defer close(ws.closed)
	defer ws.conn.Close()

	for {
		select {
		case <-ws.close:
			return
		case msg := <-ws.outgoing:
			if err := ws.writeMessage(msg); err != nil {
				// TODO: Log a write error.
				return
			}
		}
	}
}

func (ws *WebSocket) writeMessage(msg []byte) error {
	w, err := ws.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(msg)
	// TODO: Log an outgoing message.
	return err
}

func (ws *WebSocket) onClosed() {
	listeners := ws.copyListeners()
	for _, l := range listeners {
		l.onWebSocketClosed(ws)
	}
}

func (ws *WebSocket) addListener(l webSocketListener) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for _, listener := range ws.listeners {
		if l == listener {
			return
		}
	}
	ws.listeners = append(ws.listeners, l)
}

func (ws *WebSocket) removeListener(l webSocketListener) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	filtered := ws.listeners[0:0]
	for _, listener := range ws.listeners {
		if listener != l {
			filtered = append(filtered, listener)
		}
	}
	ws.listeners = filtered
}

func (ws *WebSocket) copyListeners() []webSocketListener {
	ws.mu.Lock()
	if len(ws.listeners) == 0 {
		ws.mu.Unlock()
		return nil
	}

	listeners := make([]webSocketListener, len(ws.listeners))
	copy(listeners, ws.listeners)
	ws.mu.Unlock()
	return listeners
}
