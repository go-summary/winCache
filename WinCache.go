package WinCache

import "C"
import (
	"./byteView"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 接口型函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error){
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu       sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建一个新的group实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 获得指定的group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 从缓存中获取对应的值
func (g *Group) Get(key string) (byteView.ByteView, error) {
	if key == "" {
		return byteView.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value byteView.ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (byteView.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return byteView.ByteView{}, err

	}
	value := byteView.ByteView{B: byteView.CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value byteView.ByteView) {
	g.mainCache.save(key, value)
}



