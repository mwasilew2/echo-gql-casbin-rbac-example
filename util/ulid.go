package util

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type UlidManager struct {
	entropyPool sync.Pool
}

func NewUlidManager() *UlidManager {
	return &UlidManager{
		entropyPool: sync.Pool{
			New: func() interface{} {
				entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
				return entropy
			},
		},
	}
}

func (m *UlidManager) NewULID() ulid.ULID {
	entropy := m.entropyPool.Get().(*ulid.MonotonicEntropy)
	defer m.entropyPool.Put(entropy)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
}
