package gowebgrep

import (
	"sync"
)

type Timer struct {
	Mutex        *sync.Mutex
	Milliseconds int
}

func CreateTimer(milliseconds int) *Timer {
	return &Timer{
		Mutex:        &sync.Mutex{},
		Milliseconds: milliseconds,
	}
}

// NO NEEDED !

// func (t *Timer) TriggerTimmer() {
// 	go func() {
// 		for {
// 			time.Sleep(time.Duration(t.Milliseconds))
// 		}
// 	}()
// }
