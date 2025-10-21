package main

import "fmt"

// Set 是一个泛型集合，元素类型为 T
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet 创建一个新的 Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

// Add 添加一个或多个元素到集合
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

// Remove 从集合中移除一个元素
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains 检查集合中是否包含某个元素
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Len 返回集合中元素的数量
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Clear 清空集合
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// ToSlice 将集合转换为切片（便于遍历或排序）
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}
	return result
}

// 示例使用
func main() {
	// 创建一个字符串 Set
	set := NewSet[string]()
	set.Add("apple", "banana", "orange")
	fmt.Println("初始集合:", set.ToSlice()) // [apple banana orange]（顺序可能不同）

	// 添加元素
	set.Add("grape")
	fmt.Println("添加 grape 后:", set.ToSlice())

	// 检查是否存在
	fmt.Println("包含 banana?", set.Contains("banana")) // true
	fmt.Println("包含 pear?", set.Contains("pear"))     // false

	// 删除元素
	set.Remove("banana")
	fmt.Println("删除 banana 后:", set.ToSlice())

	// 集合大小
	fmt.Println("集合大小:", set.Len()) // 3

	// 清空集合
	set.Clear()
	fmt.Println("清空后集合大小:", set.Len()) // 0
}