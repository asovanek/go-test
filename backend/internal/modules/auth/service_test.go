package auth

import (
	"context"
	"sync"
	"testing"
	"time"

	"authapp/internal/modules/user"
	"authapp/internal/platform/events"
	"authapp/internal/testdb"
	"authapp/internal/testutil"
)

func newTestService(t *testing.T) (*Service, *events.Bus) {
	t.Helper()
	db := testdb.OpenSQLite(t)
	repo := user.NewRepository(db)
	bus := events.NewBus(nil)
	cfg := testutil.TestConfig(t)
	return NewService(repo, bus, cfg), bus
}

func TestService_SignUp_SignIn(t *testing.T) {
	svc, bus := newTestService(t)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	bus.Subscribe(events.EventUserSignedUp, func(_ context.Context, e events.Event) error {
		defer wg.Done()
		evt, ok := e.(UserSignedUp)
		if !ok {
			t.Errorf("expected UserSignedUp event, got %T", e)
			return nil
		}
		if evt.Email != "new@example.com" {
			t.Errorf("event email = %q", evt.Email)
		}
		return nil
	})

	signup, err := svc.SignUp(ctx, "new@example.com", "password123")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	if signup.Token == "" {
		t.Fatal("expected token")
	}
	if signup.User.Email != "new@example.com" {
		t.Fatalf("user email = %q", signup.User.Email)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("UserSignedUp subscriber was not called")
	}

	signin, err := svc.SignIn(ctx, "new@example.com", "password123")
	if err != nil {
		t.Fatalf("signin: %v", err)
	}
	if signin.Token == "" {
		t.Fatal("expected signin token")
	}
}

func TestService_SignUp_DuplicateEmail(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	if _, err := svc.SignUp(ctx, "dup@example.com", "password123"); err != nil {
		t.Fatalf("first signup: %v", err)
	}
	_, err := svc.SignUp(ctx, "dup@example.com", "password123")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
	if err.Error() != "email already registered" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestService_SignIn_InvalidCredentials(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	if _, err := svc.SignUp(ctx, "exists@example.com", "password123"); err != nil {
		t.Fatalf("signup setup: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"unknown email", "missing@example.com", "password123"},
		{"wrong password", "exists@example.com", "wrongpassword"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.SignIn(ctx, tc.email, tc.password)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != "invalid credentials" {
				t.Fatalf("error = %q", err.Error())
			}
		})
	}
}
