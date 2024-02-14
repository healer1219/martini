package global

import "sync"

type Context struct {
	sync.Mutex
	kv map[string]any
}

func (c *Context) Set(k string, v any) {
	defer c.Unlock()
	c.Lock()
	c.kv[k] = v
}

func (c *Context) Get(k string) any {
	defer c.Unlock()
	c.Lock()
	return c.kv[k]
}

func (c *Context) Del(k string) {
	defer c.Unlock()
	c.Lock()
	delete(c.kv, k)
}

func DefaultCtx() *Context {
	return &Context{kv: make(map[string]any)}
}

func DefaultCtxWithValue(initMap map[string]any) *Context {
	return &Context{kv: initMap}
}
