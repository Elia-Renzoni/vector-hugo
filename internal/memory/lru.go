package memory

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"time"
)

type lru struct {
	hotEntries []lruEntry
	cursor     int
	poolSize   int
}

type lruEntry struct {
	data     []byte
	lastUsed time.Time
}

func newLru(poolSize int) *lru {
	return &lru{
		hotEntries: make([]lruEntry, poolSize),
		poolSize:   poolSize,
	}
}

func (l *lru) startRandomCursor(upperLimit int) {
	l.cursor = rand.IntN(upperLimit + 1)
}

func (l *lru) size() int {
	return len(l.hotEntries)
}

func (l *lru) add(entry []byte) {
	e := lruEntry{
		data:     entry,
		lastUsed: time.Now(),
	}
	l.hotEntries = append(l.hotEntries, e)

	slices.SortFunc(l.hotEntries, func(prev, next lruEntry) int {
		return prev.lastUsed.Compare(next.lastUsed)
	})
}

func (l *lru) touch(entry []byte) {
	for _, lruEntry := range l.hotEntries {
		if bytes.Equal(lruEntry.data, entry) {
			lruEntry.lastUsed = time.Now()
			break
		}
	}

	slices.SortFunc(l.hotEntries, func(prev, next lruEntry) int {
		return prev.lastUsed.Compare(next.lastUsed)
	})
}

func (l *lru) getVictim() []byte {
	return l.hotEntries[0].data
}
