package LRU

import (
	"container/list"
)

/**
FIFO 先入先出
LFU 最少使用 维护一个访问次数排序的队列
LRU 最近最少使用 如果一个数据最近访问过， 那么将来被访问的可能性就会很高
只需要维护一个队列，如果被访问过，就移动到队尾，队首肯定是最近没有访问过
*/

type Cache struct {
	maxBytes  int64                    //最大缓存大小
	nbytes    int64                    //当前缓存大小
	ll        *list.List               //当前维护的双向链表
	cache     map[string]*list.Element //缓存内容 键是字符串，值是双向链表的指针
	Onevicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

//新构造一个Cache

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		Onevicted: onEvicted,
	}
}

//获取缓存链表的长度
func (c *Cache) Len() int {
	return c.ll.Len()
}

//根据key 获取缓存中的值
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

/**
    添加元素
	1.元素存在，要进行更新，说明访问了，要放到队头
	2.元素不存在，添加进去，并且访问到了


*/

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) //value是更新的值
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(value.Len()) + int64(len(key))
	}
	//
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.ll.Remove(ele)
		c.nbytes -= int64(kv.value.Len()) + int64(len(kv.key))
		if c.Onevicted != nil {
			c.Onevicted(kv.key, kv.value)
		}
	}

}
