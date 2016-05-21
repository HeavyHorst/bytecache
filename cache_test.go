package bytecache

import (
	"testing"
	"time"
)

var result []byte

func testCache(cache Cache, t *testing.T) {
	a, _ := cache.Get("a")
	if a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	cache.Set("a", []byte("test"), -2)
	x, err := cache.Get("a")
	if err != nil {
		t.Error(err)
	}

	if string(x) != "test" {
		t.Error("x should be test is: ", string(x))
	}

	cache.Delete("a")
	x, err = cache.Get("a")
	if err != nil {
		t.Error(err)
	}
	if x != nil {
		t.Error("Found a when it should have been deleted")
	}
}

func testCacheTimes(cache Cache, t *testing.T) {
	cache.Set("a", []byte("test"), 1)
	<-time.After(2 * time.Second)
	a, err := cache.Get("a")
	if err != nil {
		t.Error(err)
	}
	if a != nil {
		t.Error("Found a when it should have been automatically deleted")
	}
}

func TestMemoryCache(t *testing.T) {
	cache := NewSimpleCache(5)
	testCache(cache, t)
}

func benchmarkWrite(cache Cache, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("test", []byte("abcdefghijklmnopqrstuvwxyz"), 1)
	}
}

func benchmarkRead(cache Cache, b *testing.B) {
	var r []byte
	cache.Set("test", []byte("abcdefghijklmnopqrstuvwxyz"), -1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, _ = cache.Get("test")
	}
	result = r
	cache.Delete("test")
}

func BenchmarkMemoryCacheWrite(b *testing.B) {
	cache := NewSimpleCache(5)
	benchmarkWrite(cache, b)
}

func BenchmarkMemoryCacheRead(b *testing.B) {
	cache := NewSimpleCache(5)
	benchmarkRead(cache, b)
}
