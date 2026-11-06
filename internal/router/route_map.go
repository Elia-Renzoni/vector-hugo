package router

import (
	"errors"
	"sync"

	"github.com/vector-hugo/internal/conf"
	"github.com/vector-hugo/internal/memory"
)

type ListDB struct {
	dbs        map[string]*memory.Arena
	mu         sync.Mutex
	memOptions conf.DBOptions
}

func NewListDB(options conf.DBOptions) *ListDB {
	return &ListDB{
		dbs:        make(map[string]*memory.Arena),
		memOptions: options,
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

	options, err := l.memOptions.Load(dsName)
	if errors.Is(err, conf.ErrDBNotFound) {
		options, err = l.memOptions.Load("standard")
	}

	arena, err := memory.NewArena(
		options.MemoryMaxBytes,
		options.MemoryThreshold,
		options.LRUPoolSize,
	)

	if err != nil {
		return nil, err
	}

	l.dbs[dsName] = arena
	return arena, nil
}
