package protectron

import (
	"sync"
	"time"
)

type registry struct {
	mu    sync.RWMutex
	users map[int64]time.Time
	stop  chan struct{}
}

func newRegistry(vacuumDelta time.Duration) *registry {
	var reg = &registry{
		users: make(map[int64]time.Time),
		stop:  make(chan struct{}),
	}
	startVacuum(vacuumDelta, reg)
	return reg
}

func (reg *registry) addUser(id int64) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.users[id] = time.Now()
}

func (reg *registry) getUserJoinTime(id int64) (time.Time, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	var joined, ok = reg.users[id]
	return joined, ok
}

func (reg *registry) stopVacuum() {
	close(reg.stop)
}

func startVacuum(delta time.Duration, reg *registry) {
	time.AfterFunc(delta, func() {
		reg.mu.Lock()
		var horizon = time.Now().Add(-delta)
		for id, joined := range reg.users {
			if horizon.After(joined) {
				delete(reg.users, id)
			}
		}
		reg.mu.Unlock()
		select {
		case <-reg.stop:
			return
		default:
			startVacuum(delta, reg)
		}
	})
}
