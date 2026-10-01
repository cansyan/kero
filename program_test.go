package kero

import (
	"sync"
	"testing"
)

func TestRequestFrameQueuesConcurrentRequests(t *testing.T) {
	p := &Program{
		frameRequests: make(chan int, 1),
	}

	const callers = 32
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			p.requestFrame(30)
		}()
	}
	wg.Wait()

	if got := len(p.frameRequests); got != 1 {
		t.Fatalf("queued requests = %d, want 1 coalesced request", got)
	}
	if p.frameC != nil || p.frameTimer != nil {
		t.Fatal("requestFrame modified timer state outside the program loop")
	}

	p.scheduleFrame(<-p.frameRequests)
	if p.frameC == nil || p.frameTimer == nil {
		t.Fatal("scheduleFrame did not arm the frame timer")
	}
	p.frameTimer.Stop()
}
