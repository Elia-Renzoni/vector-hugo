package router

import (
	"sync"

	"github.com/vector-hugo/internal/memory"
)

type ListDB struct {
	dbs map[string]*memory.Arena
	mu  sync.Mutex
}

func NewListDB() *ListDB {
	return &ListDB{
		dbs: make(map[string]*memory.Arena),
	}
}

func (l *ListDB) Route(dsName string, forceInsertion bool) (*memory.Arena, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	val, ok := l.dbs[dsName]
	if ok {
		return val, nil
	}

	if !forceInsertion {
		return nil, nil
	}

	arena, err := memory.NewArena()
	if err != nil {
		return nil, err
	}

	l.dbs[dsName] = arena
	return arena, nil
}
