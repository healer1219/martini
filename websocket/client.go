package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	clientPool = Default()
)

type PoolKey interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string
}

type ConnPool[T PoolKey] struct {
	sync.RWMutex
	ConnMap map[T]*websocket.Conn
}

func Default() *ConnPool[string] {
	return New[string]()
}

func New[T PoolKey]() *ConnPool[T] {
	return &ConnPool[T]{
		ConnMap: make(map[T]*websocket.Conn),
	}
}

func (c *ConnPool[T]) AddClient(id T, conn *websocket.Conn) error {
	defer c.Unlock()
	c.Lock()
	_, ok := c.ConnMap[id]
	if ok {
		return errors.New("conn id exist")
	}
	c.ConnMap[id] = conn
	return nil
}

func (c *ConnPool[T]) ForceAddClient(id T, conn *websocket.Conn) {
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
