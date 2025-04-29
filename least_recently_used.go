package lru

type item struct {
  value interface{}
  index int
}

type node struct {
  value interface{}
  key string
  next *node
  prev *node
}
type KV struct {
   Key string
   Value interface{}
}

type Cache interface {
  Add(string, interface{})
  Get(string) (interface{}, bool)
  Dump() []KV
} 

type ArrayCache struct {
   data map[string]*item
   recent []string
}

type LinkedListCache struct {
   data map[string]*node
   first *node
   last *node
   capacity int
   length int
}

func NewArrayCache(capacity int) Cache {
  return &ArrayCache{data: make(map[string]*item), recent: make([]string, 0, capacity)}
}

func (cache *ArrayCache) Add(key string, value interface{}) {
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

func (cache *ArrayCache) Get(key string) (interface{}, bool) {
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

func (cache *ArrayCache) Dump() []KV {
  if len(cache.recent) == 0 {
    return []KV{}
  }
  
  list := make([]KV, 0, len(cache.recent))

  for _, key := range cache.recent {
     if itm, ok := cache.data[key]; ok {
       list = append(list, KV{key, itm.value})
     } else {
       panic("expected key not found")
     }
  }
  return list
}

func NewLinkedListCache(capacity int) Cache {
  return &LinkedListCache{data: make(map[string]*node), first: nil, last: nil, capacity: capacity, length: 0 }
}

func (cache *LinkedListCache) Add(key string, value interface{}) {
  if found, ok := cache.data[key]; ok {
    found.value = value 
    cache.refresh(found)
    return
  }

  newNode := node{value: value, key: key, prev: nil, next: nil}

  if cache.length >= cache.capacity {
    if cache.first == nil || cache.last == nil {
      panic("linkedList cache with len > 0 has first or last == nil")
    }

    delete(cache.data, cache.first.key)

    if cache.capacity == 1 {
      cache.first = &newNode
      cache.last = &newNode
    } else {
      cache.first.next.prev = nil
      cache.first = cache.first.next
      cache.last.next = &newNode
      newNode.prev = cache.last
      cache.last = &newNode
    }
  } else if cache.length == 0 {
      cache.first = &newNode
      cache.last = &newNode
      cache.length = 1
  } else {
      cache.length++
      cache.last.next = &newNode
      newNode.prev = cache.last
      cache.last = &newNode
  }

  cache.data[key] = &newNode
}

func (cache *LinkedListCache) Get(key string) (interface{}, bool) {
  if found, ok := cache.data[key]; ok {
    cache.refresh(found)
    return found.value, true
  } else {
    return nil, false
  }
}

func (cache *LinkedListCache) refresh(n *node) {
  if cache.length < 2 || n == cache.last {
    return
  }
  if n == cache.first {
    cache.first.next.prev = nil
    cache.first = cache.first.next
  } else {
    n.next.prev = n.prev
    n.prev.next = n.next
  }
  n.prev = cache.last
  n.next = nil
  cache.last.next = n
  cache.last = n
}

func (cache *LinkedListCache) Dump() []KV {
  if cache.length == 0 {
    return []KV{}
  }

  list := make([]KV, 0, cache.length)
  node := cache.first
  for node.next != nil {
    list = append(list, KV{node.key, node.value})
    node = node.next
  }
  return append(list, KV{node.key, node.value})
}
