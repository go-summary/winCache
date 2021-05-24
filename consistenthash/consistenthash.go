package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 根据byte的值获得一个hash的值
type Hash func(data []byte) uint32

// Map 包括了所有hash节点值
type Map struct {
	hash     Hash
	replicas int // 虚拟节点倍数
	keys     []int // sorted //真实节点
	hashMap  map[int]string
}

// 创建一个实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 默认使用ChecksumIEEE算法，通过传入的string值生成唯一的校验码
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 增加一些真实节点(key 对应 真实节点)
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //获得该值对应的hash值
			m.keys = append(m.keys, hash) // 所有的key都要进行导入
			m.hashMap[hash] = key //虚拟节点hash所对应的真实节点
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 二分搜索进行查找
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// 删除此节点，其他的通过hash算法会逐渐去除
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}


