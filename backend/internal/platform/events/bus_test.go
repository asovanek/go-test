package events

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type testEvent struct {
	name string
}

func (e testEvent) Type() string { return e.name }

func TestBus_PublishDispatchesToSubscribers(t *testing.T) {
	bus := NewBus(nil)
	var wg sync.WaitGroup
	wg.Add(2)

	var got []string
	var mu sync.Mutex

	bus.Subscribe("test.event", func(_ context.Context, e Event) error {
		defer wg.Done()
		mu.Lock()
		got = append(got, e.Type())
		mu.Unlock()
		return nil
	})
	bus.Subscribe("test.event", func(_ context.Context, e Event) error {
		defer wg.Done()
		mu.Lock()
		got = append(got, "second-"+e.Type())
		mu.Unlock()
		return nil
	})

	bus.Publish(context.Background(), testEvent{name: "test.event"})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("subscribers were not invoked in time")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 handler calls, got %d", len(got))
	}
}

func TestBus_PublishIgnoresNilEvent(t *testing.T) {
	bus := NewBus(nil)
	called := false
	bus.Subscribe("any", func(_ context.Context, _ Event) error {
		called = true
		return nil
	})
	bus.Publish(context.Background(), nil)
	time.Sleep(50 * time.Millisecond)
	if called {
		t.Fatal("nil event should not invoke subscribers")
	}
}

func TestBus_SubscriberErrorIsLoggedNotPanicked(t *testing.T) {
	bus := NewBus(nil)
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("err.event", func(_ context.Context, _ Event) error {
		defer wg.Done()
		return errors.New("boom")
	})

	bus.Publish(context.Background(), testEvent{name: "err.event"})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("error subscriber did not run")
	}
}

func TestBus_RecoversSubscriberPanic(t *testing.T) {
	bus := NewBus(nil)
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe("panic.event", func(_ context.Context, _ Event) error {
		defer wg.Done()
		panic("subscriber panic")
	})

	bus.Publish(context.Background(), testEvent{name: "panic.event"})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("panicking subscriber did not run")
	}
}
