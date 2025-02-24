package hashmap

import (
    "fmt"
    "testing"
)

// 生成测试数据
func generateTestData(n int) (keys []string, values []int) {
    keys = make([]string, n)
    values = make([]int, n)
    for i := 0; i < n; i++ {
        keys[i] = fmt.Sprintf("key-%d", i)
        values[i] = i
    }
    return
}

// 测试普通 HashMap
func BenchmarkHashMap(b *testing.B) {
    keys, values := generateTestData(1000)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        m := NewHashMap()
        // 插入测试
        b.Run("Put", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                m.Put(keys[j], values[j])
            }
        })

        // 查找测试
        b.Run("Get", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                _, _ = m.Get(keys[j])
            }
        })

        // 删除测试
        b.Run("Remove", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                m.Remove(keys[j])
            }
        })
    }
}

// 测试 SwissHashMap
func BenchmarkSwissHashMap(b *testing.B) {
    keys, values := generateTestData(1000)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        m := NewSwissHashMap()
        // 插入测试
        b.Run("Put", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                m.Put(keys[j], values[j])
            }
        })

        // 查找测试
        b.Run("Get", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                _, _ = m.Get(keys[j])
            }
        })

        // 删除测试
        b.Run("Remove", func(b *testing.B) {
            for j := 0; j < len(keys); j++ {
                m.Remove(keys[j])
            }
        })
    }
}

// 混合操作测试
func BenchmarkMixedOperations(b *testing.B) {
    keys, values := generateTestData(1000)
    
    // 测试普通 HashMap
    b.Run("HashMap", func(b *testing.B) {
        m := NewHashMap()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            // 插入一半数据
            for j := 0; j < len(keys)/2; j++ {
                m.Put(keys[j], values[j])
            }
            // 查找操作
            for j := 0; j < len(keys)/4; j++ {
                _, _ = m.Get(keys[j])
            }
            // 删除操作
            for j := 0; j < len(keys)/4; j++ {
                m.Remove(keys[j])
            }
        }
    })

    // 测试 SwissHashMap
    b.Run("SwissHashMap", func(b *testing.B) {
        m := NewSwissHashMap()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            // 插入一半数据
            for j := 0; j < len(keys)/2; j++ {
                m.Put(keys[j], values[j])
            }
            // 查找操作
            for j := 0; j < len(keys)/4; j++ {
                _, _ = m.Get(keys[j])
            }
            // 删除操作
            for j := 0; j < len(keys)/4; j++ {
                m.Remove(keys[j])
            }
        }
    })
}