package Cache

import (
	"DisCache/ByteView"
	"DisCache/peer"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/**
每个group 都有对应一个缓存

*/

type Group struct {
	name      string
	getter    Getter //即缓存未命中时获取源数据的回调(callback)
	mainCache cache  //开始实现的并发缓存
	peers     peer.PeerPicker
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

/**
如果在缓存中，取

*/

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}

func (g *Group) RegisterPeers(peers peer.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

/**
如果不在缓存中，到数据源中取
*/

func (g *Group) load(key string) (value ByteView.ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView.ByteView{}, err
	}
	value := ByteView.ByteView{B: ByteView.CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil

}
func (g *Group) getFromPeer(peer peer.PeerGetter, key string) (ByteView.ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView.ByteView{}, err
	}
	return ByteView.ByteView{bytes}, nil
}

/**
从源数据源获取数据后也会加入到缓存中
*/
func (g *Group) populateCache(key string, value ByteView.ByteView) {
	g.mainCache.add(key, value)
}

func (g *Group) Get(key string) (ByteView.ByteView, error) {
	if key == "" {
		return ByteView.ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.GetC(key); ok {
		log.Println("GeeCache hit")
		return v, nil
	}
	return g.load(key)
}
