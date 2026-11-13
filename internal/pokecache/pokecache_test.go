package pokecache

import (
	"testing"
	"time"
)

// testInterval jest długi, żeby reapLoop nie usunął niczego
// w trakcie testów, które sprawdzają samo Add/Get.
const testInterval = 5 * time.Second

func TestAddGet(t *testing.T) {
	cache := NewCache(testInterval)

	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
	}

	// Wszystkie klucze w tym samym cache — sprawdzamy, że się nie nadpisują.
	for _, c := range cases {
		cache.Add(c.key, c.val)
	}

	for _, c := range cases {
		got, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("Get(%q): expected key to exist", c.key)
			continue
		}
		if string(got) != string(c.val) {
			t.Errorf("Get(%q) = %q, want %q", c.key, got, c.val)
		}
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(testInterval)

	got, ok := cache.Get("nie-ma-mnie")
	if ok {
		t.Errorf("Get on missing key: ok = true, want false")
	}
	if got != nil {
		t.Errorf("Get on missing key: val = %q, want nil", got)
	}
}

func TestAddReplacesValue(t *testing.T) {
	cache := NewCache(testInterval)
	key := "https://example.com"

	cache.Add(key, []byte("stara"))
	cache.Add(key, []byte("nowa"))

	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get(%q): expected key to exist", key)
	}
	if string(got) != "nowa" {
		t.Errorf("Get(%q) = %q, want %q", key, got, "nowa")
	}
}

func TestFreshEntryStaysAvailable(t *testing.T) {
	const interval = 50 * time.Millisecond
	cache := NewCache(interval)
	key := "https://example.com"

	cache.Add(key, []byte("testdata"))

	// Krócej niż interval — wpis nie może jeszcze wygasnąć.
	time.Sleep(interval / 2)

	if _, ok := cache.Get(key); !ok {
		t.Errorf("Get(%q) right after Add: expected key to exist", key)
	}
}

func TestReapLoopRemovesOldEntry(t *testing.T) {
	const interval = 10 * time.Millisecond
	cache := NewCache(interval)
	key := "https://example.com"

	cache.Add(key, []byte("testdata"))

	if _, ok := cache.Get(key); !ok {
		t.Fatalf("Get(%q) right after Add: expected key to exist", key)
	}

	// Wpis może przeżyć do następnego ticka po wygaśnięciu,
	// więc czekamy z zapasem (kilka intervali).
	time.Sleep(interval * 5)

	if _, ok := cache.Get(key); ok {
		t.Errorf("Get(%q) after %v: expected key to be reaped", key, interval*5)
	}
}
