package hashmap

import (
	"testing"
)

func TestHashMap(t *testing.T) {
	// 创建新的哈希表
	m := NewHashMap()

	// 测试 IsEmpty
	if !m.IsEmpty() {
		t.Error("新创建的哈希表应该为空")
	}

	// 测试 Put 和 Get
	m.Put("key1", "value1")
	m.Put("key2", "value2")

	// 测试 Size
	if m.Size() != 2 {
		t.Errorf("期望大小为 2，实际为 %d", m.Size())
	}

	// 测试获取存在的键
	if value, exists := m.Get("key1"); !exists || value != "value1" {
		t.Error("无法获取已存在的键值对")
	}

	// 测试获取不存在的键
	if _, exists := m.Get("nonexistent"); exists {
		t.Error("不应该找到不存在的键")
	}

	// 测试更新现有键的值
	m.Put("key1", "newvalue1")
	if value, _ := m.Get("key1"); value != "newvalue1" {
		t.Error("更新键值对失败")
	}

	// 测试删除
	if !m.Remove("key1") {
		t.Error("删除已存在的键应该返回 true")
	}
	if m.Remove("nonexistent") {
		t.Error("删除不存在的键应该返回 false")
	}

	// 测试 Clear
	m.Clear()
	if !m.IsEmpty() {
		t.Error("清空后哈希表应该为空")
	}

	// 测试大量数据插入（测试扩容）
	for i := 0; i < 100; i++ {
		m.Put(i, i*2)
	}
	if m.Size() != 100 {
		t.Errorf("期望大小为 100，实际为 %d", m.Size())
	}

	// 验证所有插入的值
	for i := 0; i < 100; i++ {
		if value, exists := m.Get(i); !exists || value != i*2 {
			t.Errorf("键 %d 的值不正确", i)
		}
	}
}

func TestHashMapWithDifferentTypes(t *testing.T) {
	m := NewHashMap()

	// 测试不同类型的键
	m.Put(42, "int key")
	m.Put(3.14, "float key")
	m.Put(true, "bool key")
	m.Put(struct{ name string }{"test"}, "struct key")

	// 验证不同类型的键
	if value, _ := m.Get(42); value != "int key" {
		t.Error("整数键测试失败")
	}
	if value, _ := m.Get(3.14); value != "float key" {
		t.Error("浮点数键测试失败")
	}
	if value, _ := m.Get(true); value != "bool key" {
		t.Error("布尔键测试失败")
	}
	if value, _ := m.Get(struct{ name string }{"test"}); value != "struct key" {
		t.Error("结构体键测试失败")
	}
}