package protectron

import (
	"sync"
	"time"
)

type registry struct {
	mu    sync.RWMutex
	users map[int]time.Time
	stop  chan struct{}
}

func newRegistry(vacuumDelta time.Duration) *registry {
	var reg = &registry{
		users: make(map[int]time.Time),
		stop:  make(chan struct{}),
	}
	startVacuum(vacuumDelta, reg)
	return reg
}

func (reg *registry) addUser(id int) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.users[id] = time.Now()
}

func (reg *registry) getUserJoinTime(id int) (time.Time, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	var joined, ok = reg.users[id]
	return joined, ok
}

func (reg *registry) userIsTooYoung(delta time.Duration, id int) bool {
	var joined, ok = reg.getUserJoinTime(id)
	if !ok {
		return false
	}
	var banHorizon = time.Now().Add(-delta)
	return joined.After(banHorizon)
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
