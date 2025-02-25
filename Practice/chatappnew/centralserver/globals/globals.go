package globals

import (
	"sync"
)

const SnapshotLength = 100
const SharedMessageBufferSize = 100

type Users struct {
	Container map[string]bool
	RWlock    sync.RWMutex
}
