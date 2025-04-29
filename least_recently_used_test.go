package lru

import (
  "testing"
  "strconv"
  "reflect"
)

func newCache(capacity int) Cache {
  return NewLinkedListCache(capacity)
}

func TestBasicAdd(t *testing.T) {
  cache := newCache(5)
  cache.Add("one", 1)
  expected := 1
  actual, ok := cache.Get("one")
  if !ok {
    t.Fatalf("expected item not retrieved")
  }
  if actual != expected {
    t.Fatalf("got: %d, expected: %d", actual, expected)
  }
}

func TestGetRefreshedItem(t *testing.T) {
  cache := newCache(2)
  cache.Add("one", 1)
  cache.Add("two", 2)
  cache.Get("one")
  cache.Add("three", 3)
  if _, ok := cache.Get("one"); !ok {
    t.Fatalf("key 'two' should not be found")
  }
  if _, ok := cache.Get("two"); ok {
    t.Fatalf("key 'two' should not be found")
  }
  if _, ok := cache.Get("three"); !ok {
    t.Fatalf("key 'three' should be found")
  }
}

func TestDump(t *testing.T) {
  cache := newCache(3)
  cache.Add("one", 1)
  cache.Add("two", 2)
  cache.Add("three", 3)

  expected := []KV {
    {"one", 1},
    {"two", 2},
    {"three", 3},
  } 
  actual := cache.Dump()
  
  if !reflect.DeepEqual(expected, actual) {
    t.Fatalf("\nexpected: %v\ngot     : %v\n", expected, actual)
  }
}

func TestDump2(t *testing.T) {
  cache := newCache(3)
  cache.Add("one", 1)
  cache.Add("two", 2)
  cache.Add("three", 3)
  cache.Get("one")
  cache.Add("four", 4)
  cache.Get("three")

  expected := []KV {
    {"one", 1},
    {"four", 4},
    {"three", 3},
  } 
  actual := cache.Dump()
  
  if !reflect.DeepEqual(expected, actual) {
    t.Fatalf("\nexpected: %v\ngot     : %v\n", expected, actual)
  }
}

func TestBasicUpdate(t *testing.T) {
  cache := newCache(5)
  cache.Add("one", 1)
  cache.Add("one", "ONE")
  expected := "ONE"
  actual, ok := cache.Get("one")
  if !ok {
    t.Fatalf("expected item not retrieved")
  }
  if actual != expected {
    t.Fatalf("got: %v, expected: %v", actual, expected)
  }
}

func TestBumpOldItem(t *testing.T) {
  cache := newCache(2)
  cache.Add("one", 1)
  cache.Add("two", 2)
  cache.Add("three", 3)

  _, ok := cache.Get("one")
  if ok {
    t.Fatalf("did not expect to find 'one'")
  }
}

// Benchmark capacity: Array vs LinkedList 
// About equal at 34K items, 
// Array is about 2x faster for smaller capacities
// LinkedList is progressively faster for higher capacities 

const capacity = 34000

func BenchmarkLruSliceImpl(b *testing.B) {
  cache := NewArrayCache(capacity)
  for n := 0; n < b.N; n++ {
    doCacheOps(cache, capacity)
  }
}

func BenchmarkLruLinkedListImpl(b *testing.B) {
  cache := NewLinkedListCache(capacity)
  for n := 0; n < b.N; n++ {
    doCacheOps(cache, capacity) 
  }
}

func doCacheOps(cache Cache, n int) {
  for i := 0; i < n; i++ {
    cache.Add(strconv.Itoa(i), i)
  }
  for i := n; i < 2*n; i++ {
    cache.Add(strconv.Itoa(i), i)
  }
  for i := n; i < 2*n; i++ {
    cache.Get(strconv.Itoa(i))
  }
  for i := n; i < 2*n; i++ {
    cache.Add(strconv.Itoa(i), i+1)
  }
}
