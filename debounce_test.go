package util

import (
    "errors"
    "sync/atomic"
    "testing"
    "time"
)

func TestDebouncer_GetValue_CachesWithinDelayAndRefetchesAfter(t *testing.T) {
    var calls int32
    val := int32(0)
    fetcher := func() (int, error) {
        atomic.AddInt32(&calls, 1)
        v := int(atomic.AddInt32(&val, 1))
        return v, nil
    }

    d := NewDebouncer[int](50*time.Millisecond, fetcher)

    // First call should fetch value 1
    got, err := d.GetValue()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got != 1 {
        t.Fatalf("expected 1, got %v", got)
    }

    // Immediate subsequent calls within delay should not refetch and should return cached 1
    got, err = d.GetValue()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got != 1 {
        t.Fatalf("expected cached 1, got %v", got)
    }

    // Sleep just under delay and check still cached
    time.Sleep(20 * time.Millisecond)
    got, _ = d.GetValue()
    if got != 1 {
        t.Fatalf("expected cached 1 before delay expiry, got %v", got)
    }

    // After delay, next call should refetch and return 2
    time.Sleep(40 * time.Millisecond)
    got, _ = d.GetValue()
    if got != 2 {
        t.Fatalf("expected refreshed 2 after delay, got %v", got)
    }

    // Verify fetcher called only twice
    if c := atomic.LoadInt32(&calls); c != 2 {
        t.Fatalf("expected fetcher calls 2, got %d", c)
    }
}

func TestDebouncer_GetValue_IgnoresErrorsAndKeepsLastValue(t *testing.T) {
    var calls int32
    current := 0
    failNext := false
    fetcher := func() (int, error) {
        atomic.AddInt32(&calls, 1)
        if failNext {
            return 0, errors.New("boom")
        }
        current++
        return current, nil
    }

    d := NewDebouncer[int](30*time.Millisecond, fetcher)

    // Initial successful fetch -> 1
    v, err := d.GetValue()
    if err != nil || v != 1 {
        t.Fatalf("expected 1,nil got %v,%v", v, err)
    }

    // Force time to pass the delay
    time.Sleep(35 * time.Millisecond)

    // Next fetch fails; Debouncer's GetValue should return lastValue and nil error per implementation
    failNext = true
    v, err = d.GetValue()
    if err != nil {
        t.Fatalf("expected nil error despite fetch failure, got %v", err)
    }
    if v != 1 {
        t.Fatalf("expected to keep last value 1 on fetch error, got %v", v)
    }

    // After another delay, a successful fetch should update value to 2
    failNext = false
    time.Sleep(35 * time.Millisecond)
    v, err = d.GetValue()
    if err != nil || v != 2 {
        t.Fatalf("expected 2,nil after recovery, got %v,%v", v, err)
    }
}

func TestPubSubDebouncer_RegisterFetchesAndPublishes(t *testing.T) {
    // Use delay >= 100ms due to constructor constraint
    delay := 120 * time.Millisecond

    var cnt int32
    fetcher := func() (int, error) {
        return int(atomic.AddInt32(&cnt, 1)), nil
    }

    d, cancel := NewPubSubDebouncer[int](delay, fetcher)
    defer cancel()

    ch := d.Register()
    defer d.Unregister(ch)

    // Expect periodic publishes roughly each delay; allow generous timeouts to avoid flakiness
    select {
    case v := <-ch:
        if v != 1 {
            t.Fatalf("expected first published value 1, got %v", v)
        }
    case <-time.After(2 * delay):
        t.Fatal("timed out waiting for first publish")
    }

    // Next cycle should publish 2
    select {
    case v := <-ch:
        if v != 2 {
            t.Fatalf("expected second published value 2, got %v", v)
        }
    case <-time.After(2 * delay):
        t.Fatal("timed out waiting for second publish")
    }
}

func TestPubSubDebouncer_GetValue_RespectsDelayAndErrors(t *testing.T) {
    delay := 120 * time.Millisecond
    var calls int32
    value := 0
    fail := false

    fetcher := func() (int, error) {
        atomic.AddInt32(&calls, 1)
        if fail {
            return 0, errors.New("fail")
        }
        value++
        return value, nil
    }

    d, cancel := NewPubSubDebouncer[int](delay, fetcher)
    defer cancel()

    // First GetValue should fetch 1
    v, err := d.GetValue()
    if err != nil || v != 1 {
        t.Fatalf("expected 1,nil got %v,%v", v, err)
    }

    // Within delay, should return cached 1 and not call fetcher
    v, err = d.GetValue()
    if err != nil || v != 1 {
        t.Fatalf("expected cached 1,nil got %v,%v", v, err)
    }

    // After delay, failure should return last value and non-nil error
    time.Sleep(delay + 10*time.Millisecond)
    fail = true
    v, err = d.GetValue()
    if err == nil {
        t.Fatalf("expected error on failed fetch")
    }
    // current implementation returns zero value alongside error on failed fetch
    if v != 0 {
        t.Fatalf("expected zero value on error, got %v", v)
    }

    // After another delay, success should update to 2
    time.Sleep(delay + 10*time.Millisecond)
    fail = false
    v, err = d.GetValue()
    if err != nil || v != 2 {
        t.Fatalf("expected 2,nil got %v,%v", v, err)
    }
}

func TestPubSubDebouncer_CancelStopsFetcher(t *testing.T) {
    delay := 120 * time.Millisecond
    var calls int32
    fetcher := func() (int, error) {
        atomic.AddInt32(&calls, 1)
        return 1, nil
    }

    d, cancel := NewPubSubDebouncer[int](delay, fetcher)
    ch := d.Register()

    // Wait for first publish
    select {
    case <-ch:
    case <-time.After(2 * delay):
        t.Fatal("timeout waiting initial publish")
    }

    // Cancel and ensure no more publishes; allow a window > delay
    cancel()
    select {
    case <-ch:
        t.Fatal("received publish after cancel")
    case <-time.After(2 * delay):
        // ok
    }

    d.Unregister(ch)
}
