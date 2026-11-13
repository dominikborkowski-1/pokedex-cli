package pokecache

import (
	"sync"
	"time"
)

// cacheEntry to pojedynczy wpis: bajty plus moment dodania,
// żeby reapLoop wiedział, które wpisy są już za stare.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache przechowuje surowe bajty pod kluczem (np. URL zapytania).
// Nie wie nic o Pokemonach ani o JSON-ie.
// Mapa jest współdzielona z goroutine reapLoop, więc każdy dostęp
// (Add, Get, reap) idzie przez ten sam mutex.
type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

// NewCache tworzy cache i startuje w tle goroutine, która co `interval`
// usuwa wpisy starsze niż `interval`.
// Zwracamy wskaźnik, bo Cache zawiera mutex — kopiowanie go byłoby błędem.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

// Add zapisuje wartość pod kluczem razem z aktualnym czasem.
// Istniejący klucz zostaje nadpisany (i jego "wiek" się resetuje).
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get zwraca zapisane bajty i true, albo nil i false, gdy klucza nie ma.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// reapLoop budzi się przy każdym ticku i sprząta stare wpisy.
// Działa przez cały czas życia programu.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.reap(time.Now(), interval)
	}
}

// reap usuwa wpisy dodane wcześniej niż now - maxAge.
// Usuwanie z mapy podczas range jest w Go bezpieczne.
func (c *Cache) reap(now time.Time, maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > maxAge {
			delete(c.entries, key)
		}
	}
}
