package sdk

import (
	"sync"
)

//ListTask 任务列表 [JobID]执行函数,并行执行时[+LogID]
type ListTask struct {
	mu   sync.RWMutex
	data map[string]*Task
}

// Set 设置数据
func (t *ListTask) Set(key string, val *Task) {
	t.mu.Lock()
	t.data[key] = val
	t.mu.Unlock()
}

// Get 获取数据
func (t *ListTask) Get(key string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[key]
}

// GetAll 获取所有数据
func (t *ListTask) GetAll() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

// GetKeys 获取keys
func (t *ListTask) GetKeys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	keys := make([]string, 0)
	for k, _ := range t.data {
		keys = append(keys, k)
	}
	return keys
}

// Del 设置数据
func (t *ListTask) Del(key string) {
	t.mu.Lock()
	delete(t.data, key)
	t.mu.Unlock()
}

// Len 长度
func (t *ListTask) Len() int {
	return len(t.data)
}

// Exists Key是否存在
func (t *ListTask) Exists(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.data[key]
	return ok
}
