package blockchain

import "sync"

// Set 用来确保添加节点是幂等的的简单方
type Set struct {
	items map[string]struct{} `json:"items"`
	sync.RWMutex
}

var itemExists = struct{}{}

// NewSet .
func NewSet() *Set {
	return &Set{
		items: make(map[string]struct{}),
	}
}

// Add 添加
func (set *Set) Add(items ...string) {
	set.Lock()
	defer set.Unlock()
	for _, item := range items {
		set.items[item] = itemExists
	}
}

// Remove 删除
func (set *Set) Remove(items ...string) {
	set.Lock()
	defer set.Unlock()
	for _, item := range items {
		delete(set.items, item)
	}
}

// Contains 是否存在
func (set *Set) Contains(item string) bool {
	set.RLock()
	defer set.RUnlock()
	if _, contains := set.items[item]; !contains {
		return false
	}
	return true
}

// Empty 是否为空
func (set *Set) Empty() bool {
	return set.Size() == 0
}

// Size 大小
func (set *Set) Size() int {
	return len(set.items)
}

// Clear 重置
func (set *Set) Clear() {
	set.Lock()
	defer set.Unlock()
	set.items = make(map[string]struct{})
}

// List 返回一个切片
func (set *Set) List() []string {
	set.RLock()
	defer set.RUnlock()
	list := make([]string, 0, 100)
	for node := range set.items {
		list = append(list, node)
	}
	return list
}
