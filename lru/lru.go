package lru
import "container/list"

// 建立cache的基础结构, 使用LRU方案进行淘汰策略
type Cache struct {
	// 可使用的最大内存
	maxBytes int64
	// 当前使用的内存
	nBytes   int64
	// 双向链表
	dl       *list.List
	// 对应map的值是所在节点的指针
	cache    map[string]*list.Element
	// 某条记录被删除时候回调的函数，可空
	OnEvicted func(key string, value Value)
}

// 建立当前节点的entry
type entry struct {
	key   string
	value Value
}
// 实现了Value接口的byteView类可以进行使用
type Value interface {
	Len() int
}

// New函数
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		dl:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}


// 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// 删除
func (c *Cache) Delete() {
	// 取出队尾的数据
	ele := c.dl.Back()
	if ele != nil {
		c.dl.Remove(ele)
		// 双向链表的值可转换为entry
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 增加或者修改
func (c *Cache) Save(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.dl.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果发现多了，需要进行淘汰
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.Delete()
	}
}

// 给cache实现len接口
func (c *Cache) Len() int{
	return c.dl.Len()
}




