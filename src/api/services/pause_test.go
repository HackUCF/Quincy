package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/HackUCF/Quincy/src/api/config"
)

// initPaused builds the standard test config with the given start-paused setting
// and runs InitServices, which is what seeds the pause state.
func initPaused(t *testing.T, startPaused bool) *config.APIConfigSpec {
	t.Helper()

	cfg := minimalServicesConfig()
	cfg.StartPaused = startPaused
	setupTeamRange(cfg)

	// resetForTest does not clear the pause state, so restore the running
	// state here or every later test in the package gets no-op checks
	t.Cleanup(resetForTest)
	t.Cleanup(func() { Unpause() })

	if err := InitServices(cfg); err != nil {
		t.Fatalf("InitServices: %v", err)
	}
	return cfg
}

func TestInitServices_startPausedSeedsState(t *testing.T) {
	for _, tc := range []struct {
		name        string
		startPaused bool
	}{
		{"starts running", false},
		{"starts paused", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initPaused(t, tc.startPaused)

			if got := IsPaused(); got != tc.startPaused {
				t.Errorf("IsPaused() = %v, want %v", got, tc.startPaused)
			}

			paused, since := PauseStatus()
			if paused != tc.startPaused {
				t.Errorf("PauseStatus() paused = %v, want %v", paused, tc.startPaused)
			}
			if since.IsZero() {
				t.Error("PauseStatus() since is zero, want the time the state was entered")
			}
		})
	}
}

func TestSetPauseState_reportsOnlyRealTransitions(t *testing.T) {
	initPaused(t, false)

	if !Pause() {
		t.Error("Pause() = false on first pause, want true")
	}
	if !IsPaused() {
		t.Error("IsPaused() = false after Pause()")
	}

	if Pause() {
		t.Error("Pause() = true when already paused, want false")
	}

	if !Unpause() {
		t.Error("Unpause() = false on first unpause, want true")
	}
	if IsPaused() {
		t.Error("IsPaused() = true after Unpause()")
	}

	if Unpause() {
		t.Error("Unpause() = true when already running, want false")
	}
}

func TestSetPauseState_advancesTimestampOnlyOnChange(t *testing.T) {
	initPaused(t, false)

	_, before := PauseStatus()

	// a no-op transition must leave the timestamp alone
	if Unpause() {
		t.Fatal("Unpause() reported a change while already running")
	}
	if _, after := PauseStatus(); !after.Equal(before) {
		t.Errorf("since moved on a no-op transition: %v -> %v", before, after)
	}

	// time.Now has coarse resolution on some platforms, so make sure the
	// clock has actually moved before asserting the timestamp advanced
	time.Sleep(time.Millisecond)

	if !Pause() {
		t.Fatal("Pause() reported no change from the running state")
	}
	paused, after := PauseStatus()
	if !paused {
		t.Error("PauseStatus() paused = false after Pause()")
	}
	if !after.After(before) {
		t.Errorf("since = %v, want a time after %v", after, before)
	}
}

func TestGetNext_pausedReturnsNoOpWithoutAdvancingQueue(t *testing.T) {
	cfg := initPaused(t, false)

	// burn a position so a failure to hold the queue is visible
	if _, err := GetNext(context.Background(), cfg, nil); err != nil {
		t.Fatalf("GetNext: %v", err)
	}
	idxBefore := servicesIdx.Load()

	if !Pause() {
		t.Fatal("Pause() reported no change from the running state")
	}

	for i := range 5 {
		svc, err := GetNext(context.Background(), cfg, nil)
		if err != nil {
			t.Fatalf("GetNext at i=%d: %v", i, err)
		}
		if !svc.NoOp {
			t.Errorf("GetNext at i=%d returned a real check while paused: %+v", i, svc)
		}
	}

	if got := servicesIdx.Load(); got != idxBefore {
		t.Errorf("servicesIdx = %d after 5 paused calls, want %d — a pause must cost no queue position", got, idxBefore)
	}

	// unpausing resumes from where the counter left off
	if !Unpause() {
		t.Fatal("Unpause() reported no change from the paused state")
	}
	svc, err := GetNext(context.Background(), cfg, nil)
	if err != nil {
		t.Fatalf("GetNext after unpause: %v", err)
	}
	if svc.NoOp {
		t.Error("GetNext returned a no-op after unpausing")
	}
}

func TestSetPauseState_exactlyOneConcurrentWinner(t *testing.T) {
	initPaused(t, false)

	const goroutines = 64

	var (
		start sync.WaitGroup
		done  sync.WaitGroup
		mu    sync.Mutex
		wins  int
	)
	start.Add(1)
	done.Add(goroutines)

	for range goroutines {
		go func() {
			defer done.Done()
			start.Wait()
			if Pause() {
				mu.Lock()
				wins++
				mu.Unlock()
			}
		}()
	}

	start.Done()
	done.Wait()

	if wins != 1 {
		t.Errorf("%d goroutines reported changing the state, want exactly 1", wins)
	}
	if !IsPaused() {
		t.Error("IsPaused() = false after a concurrent pause storm")
	}
}

func TestPauseStatus_snapshotStaysConsistentUnderLoad(t *testing.T) {
	initPaused(t, false)

	const (
		writers = 8
		readers = 8
		rounds  = 200
	)

	stop := make(chan struct{})
	var writersWG, readersWG sync.WaitGroup

	for i := range writers {
		writersWG.Add(1)
		go func() {
			defer writersWG.Done()
			for range rounds {
				if i%2 == 0 {
					Pause()
				} else {
					Unpause()
				}
			}
		}()
	}

	for range readers {
		readersWG.Add(1)
		go func() {
			defer readersWG.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				// the flag and the timestamp must come from the same
				// snapshot, so a status read can never see a zero time
				if _, since := PauseStatus(); since.IsZero() {
					t.Error("PauseStatus() returned a zero timestamp — torn snapshot")
					return
				}
				IsPaused()
			}
		}()
	}

	writersWG.Wait()
	close(stop)
	readersWG.Wait()
}
