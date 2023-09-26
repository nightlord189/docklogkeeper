package trigger

import (
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"sync"
)

type TriggersCache struct {
	data map[string][]entity.TriggerDB //container_name -> entity.TriggerDB //empty for all
	lock *sync.RWMutex
}

func NewTriggersCache() *TriggersCache {
	return &TriggersCache{
		data: make(map[string][]entity.TriggerDB, 10),
		lock: &sync.RWMutex{},
	}
}

func (c *TriggersCache) Clear() {
	c.lock.Lock()
	for key := range c.data {
		delete(c.data, key)
	}
	c.lock.Unlock()
}

func (c *TriggersCache) LoadWithAll(containerName string) []entity.TriggerDB {
	c.lock.RLock()
	specific := c.data[containerName]
	all := c.data[""]
	c.lock.RUnlock()
	return append(specific, all...)
}

func (c *TriggersCache) Load(containerName string) ([]entity.TriggerDB, bool) {
	c.lock.RLock()
	specific, ok := c.data[containerName]
	c.lock.RUnlock()
	return specific, ok
}

func (c *TriggersCache) Store(containerName string, triggers []entity.TriggerDB) {
	c.lock.Lock()
	c.data[containerName] = triggers
	c.lock.Unlock()
}
