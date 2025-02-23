package hashmap

import (
	"fmt"
	"hash/fnv"
)

// Node 表示哈希表中的节点
type Node struct {
	key   interface{}
	value interface{}
	next  *Node
}

// HashMap 实现了一个基本的哈希表
type HashMap struct {
	buckets    []*Node  // 桶数组
	size       int      // 当前元素数量
	capacity   int      // 桶数组容量
	loadFactor float64  // 负载因子阈值
}

// NewHashMap 创建一个新的哈希表实例
func NewHashMap() *HashMap {
	initialCapacity := 16
	return &HashMap{
		buckets:    make([]*Node, initialCapacity),
		size:      0,
		capacity:  initialCapacity,
		loadFactor: 0.75,
	}
}

// hash 计算键的哈希值
func (m *HashMap) hash(key interface{}) uint32 {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", key)))
	return h.Sum32()
}

// getIndex 根据哈希值获取桶索引
func (m *HashMap) getIndex(hash uint32) int {
	return int(hash) & (m.capacity - 1)
}

// Put 插入或更新键值对
func (m *HashMap) Put(key interface{}, value interface{}) {
	hash := m.hash(key)
	index := m.getIndex(hash)

	// 检查是否需要更新现有节点
	for node := m.buckets[index]; node != nil; node = node.next {
		if fmt.Sprintf("%v", node.key) == fmt.Sprintf("%v", key) {
			node.value = value
			return
		}
	}

	// 创建新节点
	newNode := &Node{
		key:   key,
		value: value,
		next:  m.buckets[index],
	}
	m.buckets[index] = newNode
	m.size++

	// 检查是否需要扩容
	if float64(m.size)/float64(m.capacity) > m.loadFactor {
		m.resize()
	}
}

// Get 获取指定键的值
func (m *HashMap) Get(key interface{}) (interface{}, bool) {
	hash := m.hash(key)
	index := m.getIndex(hash)

	for node := m.buckets[index]; node != nil; node = node.next {
		if fmt.Sprintf("%v", node.key) == fmt.Sprintf("%v", key) {
			return node.value, true
		}
	}

	return nil, false
}

// Remove 删除指定键的值
func (m *HashMap) Remove(key interface{}) bool {
	hash := m.hash(key)
	index := m.getIndex(hash)

	var prev *Node
	for node := m.buckets[index]; node != nil; node = node.next {
		if fmt.Sprintf("%v", node.key) == fmt.Sprintf("%v", key) {
			if prev == nil {
				m.buckets[index] = node.next
			} else {
				prev.next = node.next
			}
			m.size--
			return true
		}
		prev = node
	}

	return false
}

// resize 扩容哈希表
func (m *HashMap) resize() {
	oldBuckets := m.buckets
	m.capacity *= 2
	m.buckets = make([]*Node, m.capacity)
	m.size = 0

	// 重新哈希所有元素
	for _, bucket := range oldBuckets {
		for node := bucket; node != nil; node = node.next {
			m.Put(node.key, node.value)
		}
	}
}

// Size 返回哈希表中的元素数量
func (m *HashMap) Size() int {
	return m.size
}

// IsEmpty 检查哈希表是否为空
func (m *HashMap) IsEmpty() bool {
	return m.size == 0
}

// Clear 清空哈希表
func (m *HashMap) Clear() {
	m.buckets = make([]*Node, m.capacity)
	m.size = 0
}