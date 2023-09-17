package trigger

import (
	"regexp"
	"sync"
)

func clearSyncMap(syncMap *sync.Map) {
	syncMap.Range(func(key interface{}, value interface{}) bool {
		syncMap.Delete(key)
		return true
	})
}

func (a *Adapter) getRegexpFromCache(regexpStr string) *regexp.Regexp {
	if regexpStr == "" {
		return nil
	}
	gotValue, ok := a.regexpCache.Load(regexpStr)
	if !ok {
		return nil
	}
	return gotValue.(*regexp.Regexp)
}
