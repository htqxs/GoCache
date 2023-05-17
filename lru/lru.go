package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 // 允许使用的最大内存
	nbytes   int64 // 当前已使用的内存
	ll       *list.List
	cache    map[string]*list.Element // 字典，键是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

// 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache.  添加
func (c *Cache) Add(key string, value Value) {
	// 如果 cache 已经存在该 key，则是更新
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 移动到队首
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) // 更新当前所用的内存
		kv.value = value                                       // 更新值
	} else { // 不存在，则是新增
		ele := c.ll.PushFront(&entry{key, value})        // 队首添加新节点
		c.cache[key] = ele                               // 在字典中添加 key 和节点的映射关系
		c.nbytes += int64(len(key)) + int64(value.Len()) // 更新当前所用的内存
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get look ups a key's value  查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 从字典中找到对应的双向链表的节点
	if ele, ok := c.cache[key]; ok {
		// 将该节点移动队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item  删除，缓存淘汰，即移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 取到队尾节点
	if ele != nil {
		c.ll.Remove(ele) // 从链表中删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中 c.cache 删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前所用的内存
		// 如果回调函数 OnEvicted 不为 nil，则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries 用来获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
