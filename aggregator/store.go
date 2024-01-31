package main

import "github.com/DenisBytes/GoToll/types"

type MemoryStore struct {
	data map[int] float64
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(data types.Distance) error {
	m.data[data.OBUID] += data.Value
	return nil 
}