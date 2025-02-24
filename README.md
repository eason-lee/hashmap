# Go 语言实现高性能 HashMap

本项目实现了两种 HashMap 数据结构：传统的基于桶的 HashMap 和基于 Swiss Table 设计的高性能 HashMap。

## 实现方案

### 1. 传统 HashMap 实现

基于桶数组的实现，具有以下特点：
- 使用桶（bucket）存储键值对，每个桶可存储 8 个元素
- 使用 tophash 进行快速比较
- 采用溢出桶处理冲突
- 支持渐进式扩容

### 2. Swiss Table 实现

基于 Google SwissTable 设计的高性能实现：
- 使用 SIMD 指令优化查找
- 采用开放寻址法处理冲突
- 使用控制位优化查找性能
- 分组存储提升缓存命中率

## 核心特性

1. 高效的查找
   - tophash 快速比较
   - SIMD 批量操作
   - 控制位优化

2. 智能的扩容
   - 渐进式扩容
   - 负载因子自适应
   - 双 buffer 设计

3. 内存优化
   - 紧凑的内存布局
   - 高效的空间利用
   - 缓存友好的访问模式

## 性能特征

### 传统 HashMap
- 查找：平均 O(1)
- 插入：平均 O(1)
- 删除：平均 O(1)
- 适合一般使用场景

### Swiss Table HashMap
- 查找：接近 O(1)，SIMD 加速
- 插入：接近 O(1)，开放寻址
- 删除：接近 O(1)，标记删除
- 适合高性能场景

## 使用示例

```go
// 传统 HashMap
map1 := NewHashMap()
map1.Put("key1", "value1")
value, exists := map1.Get("key1")

// Swiss Table HashMap
map2 := NewSwissHashMap()
map2.Put("key1", "value1")
value, exists = map2.Get("key1")
```

## 性能对比
基准测试显示：

- Swiss Table 在查找操作上性能更优
- 传统 HashMap 在混合操作场景下更稳定
- Swiss Table 在大数据量场景下优势明显
## 后续优化方向
1. 并发安全
   
   - 分段锁实现
   - 无锁算法研究
2. SIMD 优化
   
   - 支持更多 SIMD 指令
   - 优化指令序列
3. 内存管理
   
   - 内存池优化
   - GC 友好设计
4. 功能扩展
   
   - 迭代器支持
   - 序列化接口
   - 事件回调