package router

import "github.com/vector-hugo/internal/memory"

type ListDB map[string]*memory.Arena

func (l ListDB) Route(dsName string, forceInsertion bool) (*memory.Arena, error) {
	val, ok := l[dsName]
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

	l[dsName] = arena
	return arena, nil
}
