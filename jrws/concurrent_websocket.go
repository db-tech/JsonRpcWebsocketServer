package jrws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ConcurrentWebsocket struct {
	Ws     *websocket.Conn
	RWLock sync.RWMutex
}

func (cw *ConcurrentWebsocket) WriteJSON(v interface{}) error {
	cw.RWLock.Lock()
	defer cw.RWLock.Unlock()
	return cw.Ws.WriteJSON(v)
}

func (cw *ConcurrentWebsocket) ReadJSON(v interface{}) error {
	return cw.Ws.ReadJSON(v)
}
