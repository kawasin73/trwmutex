package trwmutex

import (
	"testing"
	"time"
)

func TestTRWMutex(t *testing.T) {
	var mu TRWMutex

	mu.Lock()

	if mu.TryLock() {
		t.Error("unexpectedly success to try lock")
	}
	if mu.TryRLock() {
		t.Error("unexpectedly success to try rlock")
	}

	mu.Unlock()

	if !mu.TryLock() {
		t.Error("unexpectedly failed to try lock")
	}

	mu.Unlock()

	mu.RLock()

	// RLock is shared lock
	mu.RLock()

	if mu.TryLock() {
		t.Error("unexpectedly success to try lock to read locked mutex")
	}
	if !mu.TryRLock() {
		t.Error("unexpectedly failed to try rlock")
	}

	mu.RUnlock()
	mu.RUnlock()
	mu.RUnlock()
}

func TestTRWMutex_LockOrder(t *testing.T) {
	t.Run("Lock -> Lock -> RLock", func(t *testing.T) {
		var mu TRWMutex

		// Lock -> Lock -> RLock
		mu.Lock()

		ch1 := make(chan struct{})
		go func() {
			mu.Lock()
			mu.Unlock()
			close(ch1)
		}()

		ch2 := make(chan struct{})
		go func() {
			mu.RLock()
			mu.RUnlock()
			close(ch2)
		}()

		select {
		case <-ch1:
			t.Error("unexpectedly success two write lock")
		case <-time.After(100 * time.Millisecond):
		}

		select {
		case <-ch2:
			t.Error("unexpectedly success read lock")
		case <-time.After(100 * time.Millisecond):
		}

		mu.Unlock()

		// the order of RLock and Lock is not deterministic
		select {
		case <-ch1:
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpectedly failed to lock")
		}

		select {
		case <-ch2:
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpectedly failed to read lock")
		}
	})

	t.Run("RLock -> Lock -> RLock", func(t *testing.T) {
		var mu TRWMutex

		mu.RLock()

		ch1 := make(chan struct{})
		go func() {
			mu.Lock()
			close(ch1)
		}()

		// wait to ensure following RLock() comes after Lock().
		select {
		case <-ch1:
			t.Error("unexpectedly success two write lock")
		case <-time.After(100 * time.Millisecond):
		}

		ch2 := make(chan struct{})
		go func() {
			mu.RLock()
			close(ch2)
		}()

		select {
		case <-ch1:
			t.Error("unexpectedly success two write lock")
		case <-time.After(100 * time.Millisecond):
		}

		// RLock() after Lock() always blocks
		select {
		case <-ch2:
			t.Error("unexpectedly success read lock")
		case <-time.After(100 * time.Millisecond):
		}

		mu.RUnlock()

		// Lock() success first
		select {
		case <-ch1:
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpectedly failed to lock")
		}

		mu.Unlock()

		// RLock() success later
		select {
		case <-ch2:
		case <-time.After(100 * time.Millisecond):
			t.Error("unexpectedly failed to read lock")
		}

		mu.RUnlock()
	})
}

func TestTRWMutex_TryNeverBlock(t *testing.T) {
	t.Run("Lock -> Lock -> RLock", func(t *testing.T) {
		var mu TRWMutex

		// Lock -> Lock -> RLock
		mu.Lock()

		go func() {
			mu.Lock()
			mu.Unlock()
		}()

		go func() {
			mu.RLock()
			mu.RUnlock()
		}()

		// wait to ensure Lock() and RLock() are executed
		time.Sleep(100 * time.Millisecond)

		if mu.TryLock() {
			t.Error("unexpectedly success to trylock")
		}
		if mu.TryRLock() {
			t.Error("unexpectedly success to tryrlock")
		}
	})

	t.Run("RLock -> Lock -> RLock", func(t *testing.T) {
		var mu TRWMutex

		// Lock -> Lock -> RLock
		mu.RLock()

		go func() {
			mu.Lock()
			mu.Unlock()
		}()

		// wait to ensure Lock()is executed
		time.Sleep(100 * time.Millisecond)

		go func() {
			mu.RLock()
			mu.RUnlock()
		}()

		// wait to ensure RLock() is executed
		time.Sleep(100 * time.Millisecond)

		if mu.TryLock() {
			t.Error("unexpectedly success to trylock")
		}
		if mu.TryRLock() {
			t.Error("unexpectedly success to tryrlock")
		}
	})
}
