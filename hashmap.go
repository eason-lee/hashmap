package hashmap

import (
	"fmt"
	"hash/fnv"
)

// 新增 bucket 结构
type bucket struct {
    tophash  [8]uint8     // 存储 hash 值的高 8 位，用于快速比较
    keys     [8]interface{}
    values   [8]interface{}
    overflow *bucket      // 溢出桶指针
}

// HashMap 实现了一个基本的哈希表
type HashMap struct {
    buckets    []*bucket  // 桶数组
    size       int        // 当前元素数量
    capacity   int        // 桶数组容量
    loadFactor float64    // 负载因子阈值
}

// NewHashMap 创建一个新的哈希表实例
func NewHashMap() *HashMap {
    initialCapacity := 16
    return &HashMap{
        buckets:    make([]*bucket, initialCapacity),
        size:       0,
        capacity:   initialCapacity,
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

// 获取 hash 值的高 8 位
func tophash(hash uint32) uint8 {
    return uint8(hash >> 24)
}

// Put 插入或更新键值对
func (m *HashMap) Put(key interface{}, value interface{}) {
    hash := m.hash(key)
    index := m.getIndex(hash)
    top := tophash(hash)

    // 如果桶不存在，创建新桶
    if m.buckets[index] == nil {
        m.buckets[index] = &bucket{}
    }

    // 遍历桶及其溢出桶
    for b := m.buckets[index]; b != nil; b = b.overflow {
        // 查找空位或已存在的键
        for i := 0; i < 8; i++ {
            if b.tophash[i] == 0 { // 找到空位
                b.tophash[i] = top
                b.keys[i] = key
                b.values[i] = value
                m.size++
                return
            }
            if b.tophash[i] == top && fmt.Sprintf("%v", b.keys[i]) == fmt.Sprintf("%v", key) {
                b.values[i] = value // 更新已存在的键
                return
            }
        }
        
        // 当前桶已满，需要创建或使用溢出桶
        if b.overflow == nil {
            b.overflow = &bucket{}
        }
    }

    // 检查是否需要扩容
    if float64(m.size)/float64(m.capacity*8) > m.loadFactor {
        m.resize()
    }
}

// Get 获取指定键的值
func (m *HashMap) Get(key interface{}) (interface{}, bool) {
    hash := m.hash(key)
    index := m.getIndex(hash)
    top := tophash(hash)

    // 遍历桶及其溢出桶
    for b := m.buckets[index]; b != nil; b = b.overflow {
        for i := 0; i < 8; i++ {
            if b.tophash[i] == top && fmt.Sprintf("%v", b.keys[i]) == fmt.Sprintf("%v", key) {
                return b.values[i], true
            }
        }
    }

    return nil, false
}

// Remove 删除指定键的值
func (m *HashMap) Remove(key interface{}) bool {
    hash := m.hash(key)
    index := m.getIndex(hash)
    top := tophash(hash)

    // 遍历桶及其溢出桶
    for b := m.buckets[index]; b != nil; b = b.overflow {
        for i := 0; i < 8; i++ {
            if b.tophash[i] == top && fmt.Sprintf("%v", b.keys[i]) == fmt.Sprintf("%v", key) {
                // 清空该位置
                b.tophash[i] = 0
                b.keys[i] = nil
                b.values[i] = nil
                m.size--
                return true
            }
        }
    }
    return false
}

// resize 扩容哈希表
func (m *HashMap) resize() {
    oldBuckets := m.buckets
    oldCapacity := m.capacity
    m.capacity *= 2
    m.buckets = make([]*bucket, m.capacity)
    m.size = 0

    // 重新哈希所有元素
    for i := 0; i < oldCapacity; i++ {
        for b := oldBuckets[i]; b != nil; b = b.overflow {
            for j := 0; j < 8; j++ {
                if b.tophash[j] != 0 {
                    m.Put(b.keys[j], b.values[j])
                }
            }
        }
    }
}

// Clear 清空哈希表
func (m *HashMap) Clear() {
    m.buckets = make([]*bucket, m.capacity)
    m.size = 0
}

// Size 和 IsEmpty 方法保持不变，因为它们只依赖于 size 字段
func (m *HashMap) Size() int {
	return m.size
}

// IsEmpty 检查哈希表是否为空
func (m *HashMap) IsEmpty() bool {
	return m.size == 0
}
