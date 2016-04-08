package picasa;

import (
	"time"
	"sync"
);

type cacheData struct {
	mtime	int64
	data	interface{}
}

type cacheMap struct {
	sync.RWMutex
	expire int64
	data map[string]*cacheData
}

func newCacheMap() cacheMap {
	return cacheMap {
		expire: CACHE_EXPIRE,
		data: make(map[string]*cacheData),
	}
}

func (obj *cacheMap) SetExpire(seconds int64) {
	obj.expire = seconds
}

func (obj *cacheMap) GetEntry(key string) interface{} {
	var time = time.Now().Unix()

	obj.RLock()
	var keyData = obj.data[key]
	obj.RUnlock()

	if (keyData == nil || (time > (keyData.mtime + obj.expire))) {
		return nil
	} else {
		return keyData.data
	}
}

func (obj *cacheMap) SetEntry(key string, data interface{}) {
	var time = time.Now().Unix()

	var keyData = &cacheData {
		mtime: time,
		data: data,
	}

	obj.RLock()
	obj.data[key] = keyData
	obj.RUnlock()
	return
}

func (obj *cacheMap) DelEntry(key string) {
	obj.RLock()
	delete(obj.data, key)
	obj.RUnlock()
	return
}
