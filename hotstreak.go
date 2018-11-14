package hotstreak

import (
	"sync"
	"time"
)

// Hotstreak is the main structure of the library
type Hotstreak struct {
	active     bool
	hot        bool
	mux        *sync.Mutex
	counter    int
	Limit      int
	HotWait    time.Duration
	ActiveWait time.Duration
	notifier   chan uint8
}

// Config is the structure of the configuration that can be injected to the lib
type Config struct {
	Limit      int           // Describes how many times we have to hit before a streak becomes hot
	HotWait    time.Duration // Describes the amount of time we are waiting before declaring a cool down
	ActiveWait time.Duration // Describes the amount of time we are waiting to check on a streak being active
}

var (
	// DEACTIVATED - Status sign for the streak getting deactivated
	DEACTIVATED uint8
	// ACTIVATED - Status sign for the streak getting activated
	ACTIVATED uint8 = 1
)

// New creates a new instance of Hotstreak
func New(config Config) *Hotstreak {
	limit := config.Limit
	if limit == 0 {
		limit = 20
	}
	hotwait := config.HotWait
	if hotwait == 0 {
		hotwait = time.Minute * 5
	}
	activeWait := config.ActiveWait
	if activeWait == 0 {
		activeWait = time.Minute * 5
	}
	return &Hotstreak{Limit: limit, HotWait: hotwait, ActiveWait: activeWait, mux: &sync.Mutex{}, notifier: make(chan uint8)}
}

func (hs *Hotstreak) coolDown() {
	if hs == nil {
		return
	}
	time.Sleep(hs.HotWait)
	hs.mux.Lock()
	defer hs.mux.Unlock()
	hs.hot = false
	hs.counter = 0
}

func (hs *Hotstreak) dieSlowly() {
	if hs == nil {
		return
	}
	select {
	case <-hs.notifier:
		return
	case <-time.After(hs.ActiveWait):
		hs.mux.Lock()
		defer hs.mux.Unlock()
		if hs.hot || hs.counter > 0 {
			go hs.dieSlowly()
			hs.counter = 0
			return
		}
		hs.active = false
	}
}

// Hit is to ping the service increasing it's hotness
func (hs *Hotstreak) Hit() *Hotstreak {
	if hs == nil {
		return nil
	}
	hs.mux.Lock()
	defer hs.mux.Unlock()
	if hs.hot {
		return hs
	}
	hs.counter++
	if hs.counter >= hs.Limit {
		hs.hot = true
		go hs.coolDown()
	}
	return hs
}

// Activate turns on the streak
func (hs *Hotstreak) Activate() *Hotstreak {
	if hs == nil {
		return nil
	}
	if hs.active {
		hs.notifier <- ACTIVATED
	}
	hs.active = true
	go hs.dieSlowly()
	return hs
}

// Deactivate turns off the streak
func (hs *Hotstreak) Deactivate() *Hotstreak {
	if hs == nil {
		return nil
	}
	if hs.active {
		hs.notifier <- 0
	}
	hs.active = false
	hs.hot = false
	hs.counter = 0
	return hs
}

// IsHot is for getting the hot status of Hotstreak
func (hs *Hotstreak) IsHot() bool {
	if hs == nil {
		return false
	}
	return hs.hot
}

// IsActive is for getting if the streak is active at all or not
func (hs *Hotstreak) IsActive() bool {
	if hs == nil {
		return false
	}
	return hs.active
}
