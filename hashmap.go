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
    buckets      []*bucket  // 桶数组
    oldbuckets   []*bucket  // 扩容时的旧桶数组
    size         int        // 当前元素数量
    capacity     int        // 桶数组容量
    loadFactor   float64    // 负载因子阈值
    resizing     bool       // 是否正在扩容
    resizeIndex  int        // 扩容进度
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

// resize 开始扩容过程
func (m *HashMap) resize() {
    if m.resizing {
        return // 已经在扩容中
    }
    m.oldbuckets = m.buckets
    m.buckets = make([]*bucket, m.capacity*2)
    m.resizing = true
    m.resizeIndex = 0
}

// evacuate 迁移一个桶
func (m *HashMap) evacuate(index int) {
    if !m.resizing || index >= len(m.oldbuckets) {
        return
    }

    oldBucket := m.oldbuckets[index]
    if oldBucket == nil {
        m.resizeIndex++
        return
    }

    // 迁移当前桶及其溢出桶中的所有元素
    for b := oldBucket; b != nil; b = b.overflow {
        for i := 0; i < 8; i++ {
            if b.tophash[i] != 0 {
                // 重新计算在新桶数组中的位置
                hash := m.hash(b.keys[i])
                newIndex := m.getIndex(hash)
                
                // 放入新桶
                if m.buckets[newIndex] == nil {
                    m.buckets[newIndex] = &bucket{}
                }
                m.putToNewBucket(b.keys[i], b.values[i], hash)
            }
        }
    }

    // 清除旧桶
    m.oldbuckets[index] = nil
    m.resizeIndex++

    // 检查是否完成扩容
    if m.resizeIndex >= len(m.oldbuckets) {
        m.oldbuckets = nil
        m.resizing = false
        m.capacity *= 2
    }
}

// putToNewBucket 将元素放入新桶（仅在扩容时使用）
func (m *HashMap) putToNewBucket(key interface{}, value interface{}, hash uint32) {
    index := m.getIndex(hash)
    top := tophash(hash)

    // 确保目标桶存在
    if m.buckets[index] == nil {
        m.buckets[index] = &bucket{}
    }

    // 遍历桶及其溢出桶
    for b := m.buckets[index]; ; b = b.overflow {
        // 查找空位
        for i := 0; i < 8; i++ {
            if b.tophash[i] == 0 {
                b.tophash[i] = top
                b.keys[i] = key
                b.values[i] = value
                return
            }
        }
        
        // 当前桶已满，创建或使用溢出桶
        if b.overflow == nil {
            b.overflow = &bucket{}
        }
        b = b.overflow
    }
}

// Put 插入或更新键值对
func (m *HashMap) Put(key interface{}, value interface{}) {
    // 如果正在扩容，迁移一个桶
    if m.resizing {
        m.evacuate(m.resizeIndex)
    }

    hash := m.hash(key)
    index := m.getIndex(hash)
    top := tophash(hash)

    // 如果正在扩容，需要同时检查旧桶
    if m.resizing {
        oldIndex := index & (m.capacity/2 - 1)
        if oldBucket := m.oldbuckets[oldIndex]; oldBucket != nil {
            // 在旧桶中查找并更新
            for b := oldBucket; b != nil; b = b.overflow {
                for i := 0; i < 8; i++ {
                    if b.tophash[i] == top && fmt.Sprintf("%v", b.keys[i]) == fmt.Sprintf("%v", key) {
                        b.values[i] = value
                        return
                    }
                }
            }
        }
    }

    // 原有的插入逻辑
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

    // 如果正在扩容，先检查旧桶
    if m.resizing {
        oldIndex := index & (m.capacity/2 - 1)
        if oldBucket := m.oldbuckets[oldIndex]; oldBucket != nil {
            for b := oldBucket; b != nil; b = b.overflow {
                for i := 0; i < 8; i++ {
                    if b.tophash[i] == top && fmt.Sprintf("%v", b.keys[i]) == fmt.Sprintf("%v", key) {
                        return b.values[i], true
                    }
                }
            }
        }
    }

    // 检查新桶
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

    // 如果正在扩容，先检查旧桶
    if m.resizing {
        oldIndex := index & (m.capacity/2 - 1)
        if oldBucket := m.oldbuckets[oldIndex]; oldBucket != nil {
            for b := oldBucket; b != nil; b = b.overflow {
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
        }
    }

    // 检查新桶
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
