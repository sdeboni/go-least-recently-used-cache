package lru

func TestMe() string {
  return "hello world"
}

type item struct {
  value interface{}
  index int
}

type Cache struct {
   data map[string]*item
   recent []string
}

func NewCache(capacity int) *Cache {
  return &Cache{data: make(map[string]*item), recent: make([]string, 0, capacity)}
}

func (cache *Cache) Add(key string, value interface{}) {
  if _, ok := cache.Get(key); ok {
    cache.data[key].value = value
  } else {
    if len(cache.recent) == cap(cache.recent) {
      oldestKey := cache.recent[0]
      delete(cache.data, oldestKey)
      cache.recent = cache.recent[1:]
    }
    cache.recent = append(cache.recent, key)
    itm := item{value: value, index: len(cache.recent)-1}
    cache.data[key] = &itm
  }
} 

func (cache *Cache) Get(key string) (interface{}, bool) {
  if found, ok := cache.data[key]; ok {
    if found.index == len(cache.recent) - 1 {
       return found.value, true
    }

    if found.index == 0 {
      cache.recent = cache.recent[1:]
    } else if found.index < len(cache.recent) -1 {
      cache.recent = append(cache.recent[0:found.index], cache.recent[found.index + 1:]...)
    }
    cache.recent = append(cache.recent, key)
    found.index = len(cache.recent) - 1

    return found.value, true
  } else {
    return nil, false
  }
}
