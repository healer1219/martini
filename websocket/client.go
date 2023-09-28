package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	clientPool = ConnPool{ConnMap: make(map[string]*websocket.Conn)}
)

type ConnPool struct {
	sync.RWMutex
	ConnMap map[string]*websocket.Conn
}

func (c *ConnPool) New() *ConnPool {
	return &ConnPool{
		ConnMap: make(map[string]*websocket.Conn),
	}
}

func (c *ConnPool) AddClient(id string, conn *websocket.Conn) error {
	defer c.Unlock()
	c.Lock()
	_, ok := c.ConnMap[id]
	if ok {
		return errors.New("id exist")
	}
	c.ConnMap[id] = conn
	return nil
}

func (c *ConnPool) ForceAddClient(id string, conn *websocket.Conn) {
	defer c.Unlock()
	c.Lock()
	c.ConnMap[id] = conn
}

func AddClient(id string, conn *websocket.Conn) error {
	return clientPool.AddClient(id, conn)
}

func ForceAddClient(id string, conn *websocket.Conn) {
	clientPool.ForceAddClient(id, conn)
}
