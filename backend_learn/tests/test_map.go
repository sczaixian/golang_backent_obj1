package tests

import (
	"fmt"
	"sync"
)

func TestMap() {
	/*
	*  key 可以是任何可以比较的类型，如：整数、字符串等
	*  标准的 map 不是线程安全， 需要用 mutex、rwmutex保证线程安全 使用sync.map
	*  使用没有初始化的 map 值为 nil， 直接赋值会导致运行时错误
	 */

	var myMap map[string]int
	fmt.Println(myMap)
	myMap = make(map[string]int)
	fmt.Println(myMap)
	myMap = map[string]int{
		"a": 1,
		"b": 2,
	}
	fmt.Println(myMap)

	myMap["c"] = 3
	delete(myMap, "a")
	myMap["b"] = 4

	val, exists := myMap["b"]
	// val := myMap["b"]
	fmt.Println(val, exists)

	// Map的遍历顺序是随机的，每次遍历的顺序可能都不一样
	for k, v := range myMap {
		fmt.Println(k, v)
	}

	// https://cloud.tencent.com/developer/article/2400014

}

func TestSyncMap() { // 主要用于读多写少的场景
	/*
	 *  无需初始化，不需要像内置map那样使用make函数初始化，可以直接声明后使用
	 *  Load, Store, LoadOrStore, Delete, 和Range

	 *  内置map加锁的方式效率不高
	 *  key 的集合基本不变，但是Value会并发更新：在这种场景下，sync.Map通过将热点数据分离出来，减少了锁的争用，提高了性能。
	 *  Key-Value对的添加和删除操作比较少，但是读操作非常频繁：sync.Map在读取操作上做了优化，读操作通常无需加锁，这大大提高了并发读的性能。

	 *  使用只读和读写分离的数据结构，减少了锁的争用
	 *  标记删除而不是立即删除来提高性能
	 *  根据实际的使用模式动态调整内部数据结构，以优化性能

	 * 并发环境下的缓存系统：缓存项被频繁读取，但更新和删除操作较少。
	 * 长时间运行的监听器列表：监听器被添加后很少改变，但可能会被频繁触发。
	 * 全局状态和配置：全局配置可能会在程序启动时被设置，之后只会被读取。
	 */

	var m sync.Map // 也可以使用 m := new(sync.Map) 来创建一个sync.Map实例
	value, ok := m.Load(key)
	if ok {
		// key存在，value是对应的值
	} else {
		// key不存在
	}

	m.Store(key, value) // 如果值存在会覆盖

	actual, loaded := m.LoadOrStore(key, value)
	if loaded {
		// 键已经存在，actual是已存在的值
	} else {
		// 键不存在，已存储提供的值，actual是提供的值
	}

	m.Delete(key)

	// Range方法不保证每次迭代的顺序
	// 在迭代过程中如果有其他goroutine修改map，迭代器可能会反映这些修改。
	m.Range(func(key, value interface{}) bool {
		// 使用key和value
		// 如果要停止迭代，返回false
		return true
	})

}
