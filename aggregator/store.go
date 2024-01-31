package main

import "github.com/DenisBytes/GoToll/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		
	}
}

func (m *MemoryStore) Insert(data types.Distance) error {
	return nil
}