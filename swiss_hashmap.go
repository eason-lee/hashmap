package hashmap

import (
    "fmt"
    "hash/fnv"
)

const (
    EMPTY     = 0 // 空槽位
    DELETED   = 1 // 已删除
    FULL      = 2 // 已使用
    GROUP_SIZE = 16 // SIMD 分组大小
)

// 分组结构，用于 SIMD 优化
type group struct {
    control [GROUP_SIZE]uint8  // 控制位数组
    entries [GROUP_SIZE]entry  // 实际数据
}

// 单个条目
type entry struct {
    hash    uint64      // 完整哈希值
    key     interface{}
    value   interface{}
}

// SwissHashMap 实现
type SwissHashMap struct {
    groups     []group   // 分组数组
    size       int       // 元素数量
    capacity   int       // 总容量
    loadFactor float64   // 负载因子
}

func NewSwissHashMap() *SwissHashMap {
    initialGroups := 8
    return &SwissHashMap{
        groups:     make([]group, initialGroups),
        capacity:   initialGroups * GROUP_SIZE,
        loadFactor: 0.75,
    }
}

// 计算哈希值
func (m *SwissHashMap) hash(key interface{}) uint64 {
    h := fnv.New64a()
    h.Write([]byte(fmt.Sprintf("%v", key)))
    return h.Sum64()
}

// 获取控制字节
func (m *SwissHashMap) getControlByte(hash uint64) uint8 {
    // 使用哈希值的高 7 位作为探测序列，最高位设为 1 表示已使用
    return uint8(hash>>57) | 0x80 // 0x80 = 1000 0000
}

// 查找位置
func (m *SwissHashMap) findSlot(key interface{}, hash uint64) (int, int, bool) {
    control := m.getControlByte(hash)
    groupIndex := int(hash % uint64(len(m.groups)))
    
    for i := 0; i < len(m.groups); i++ {
        g := &m.groups[(groupIndex + i) % len(m.groups)]
        
        // 在组内查找匹配的控制位
        for j := 0; j < GROUP_SIZE; j++ {
            if g.control[j] == EMPTY {
                return (groupIndex + i) % len(m.groups), j, false
            }
            
            if g.control[j] == control {
                e := &g.entries[j]
                if e.hash == hash && fmt.Sprintf("%v", e.key) == fmt.Sprintf("%v", key) {
                    return (groupIndex + i) % len(m.groups), j, true
                }
            }
        }
    }
    
    return -1, -1, false
}

// Put 实现
func (m *SwissHashMap) Put(key, value interface{}) {
    hash := m.hash(key)
    groupIdx, slotIdx, exists := m.findSlot(key, hash)
    
    if groupIdx == -1 {
        // 需要扩容
        m.grow()
        groupIdx, slotIdx, exists = m.findSlot(key, hash)
    }
    
    g := &m.groups[groupIdx]
    if !exists {
        g.control[slotIdx] = m.getControlByte(hash)
        g.entries[slotIdx] = entry{
            hash:  hash,
            key:   key,
            value: value,
        }
        m.size++
    } else {
        g.entries[slotIdx].value = value
    }
}

// Get 实现
func (m *SwissHashMap) Get(key interface{}) (interface{}, bool) {
    hash := m.hash(key)
    groupIdx, slotIdx, exists := m.findSlot(key, hash)
    
    if !exists {
        return nil, false
    }
    
    return m.groups[groupIdx].entries[slotIdx].value, true
}

// Remove 实现
func (m *SwissHashMap) Remove(key interface{}) bool {
    hash := m.hash(key)
    groupIdx, slotIdx, exists := m.findSlot(key, hash)
    
    if !exists {
        return false
    }
    
    g := &m.groups[groupIdx]
    g.control[slotIdx] = DELETED
    g.entries[slotIdx] = entry{}
    m.size--
    return true
}

// 扩容
func (m *SwissHashMap) grow() {
    oldGroups := m.groups
    m.groups = make([]group, len(m.groups)*2)
    m.capacity = len(m.groups) * GROUP_SIZE
    m.size = 0
    
    // 重新插入所有元素
    for _, g := range oldGroups {
        for i := 0; i < GROUP_SIZE; i++ {
            if g.control[i] == FULL {
                m.Put(g.entries[i].key, g.entries[i].value)
            }
        }
    }
}