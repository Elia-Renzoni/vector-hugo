package memory

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"time"
)

type Lru struct {
	hotEntries []*lruEntry
	poolSize   int
	cursor     int
}

type lruEntry struct {
	data     []byte
	lastUsed time.Time
}

func NewLru(poolSize int) *Lru {
	return &Lru{
		hotEntries: make([]*lruEntry, 0, poolSize),
		poolSize:   poolSize,
	}
}

func (l *Lru) loadRandomCursor(upperLimit int) {
	l.cursor = rand.IntN(upperLimit + 1)
}

func (l *Lru) Size() int {
	return len(l.hotEntries)
}

func (l *Lru) Add(entry []byte) {
	e := lruEntry{
		data:     entry,
		lastUsed: time.Now(),
	}

	l.hotEntries = append(l.hotEntries, &e)
}

func (l *Lru) Touch(entry []byte) {
	for _, lruEntry := range l.hotEntries {
		if bytes.Equal(lruEntry.data, entry) {
			lruEntry.lastUsed = time.Now()
			break
		}
	}

	slices.SortStableFunc(l.hotEntries, func(prev, next *lruEntry) int {
		return prev.lastUsed.Compare(next.lastUsed)
	})
}

func (l *Lru) GetVictim() []byte {
	if len(l.hotEntries) == 0 {
		return nil
	}
	return l.hotEntries[0].data
}
