package lru

import "testing"

func TestBasicAdd(t *testing.T) {
  cache := NewCache(5)
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

func TestBasicUpdate(t *testing.T) {
  cache := NewCache(5)
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
  cache := NewCache(2)
  cache.Add("one", 1)
  cache.Add("two", 2)
  cache.Add("three", 3)

  _, ok := cache.Get("one")
  if ok {
    t.Fatalf("did not expect to find 'one'")
  }
}



