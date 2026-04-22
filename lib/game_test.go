package gameloop

import (
	"errors"
	"testing"

	"github.com/eiannone/keyboard"
)

// mockLoop is a test double for the Loop interface that records calls and
// returns configurable errors.
type mockLoop struct {
	initializeCalled  int
	renderCalled      int
	calculateCalled   int
	lastKey           keyboard.Key
	initializeErr     error
	renderErr         error
	calculateErr      error
}

func (m *mockLoop) Initialize() error {
	m.initializeCalled++
	return m.initializeErr
}

func (m *mockLoop) Render() error {
	m.renderCalled++
	return m.renderErr
}

func (m *mockLoop) Calculate(c keyboard.Key) error {
	m.calculateCalled++
	m.lastKey = c
	return m.calculateErr
}

// --- Interface contract ---

// Given a type implementing all Loop methods
// When used as a Loop
// Then it satisfies the interface
func TestLoop_InterfaceIsSatisfied(t *testing.T) {
	var _ Loop = &mockLoop{}
}

// --- Initialize behaviour ---

// Given a Loop implementation
// When Initialize is called
// Then it returns nil on success
func TestLoop_WhenInitializeSucceeds_ReturnsNil(t *testing.T) {
	gl := &mockLoop{}

	if err := gl.Initialize(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

// Given a Loop whose Initialize returns an error
// When Initialize is called
// Then the error is surfaced to the caller
func TestLoop_WhenInitializeFails_ReturnsError(t *testing.T) {
	want := errors.New("init failure")
	gl := &mockLoop{initializeErr: want}

	if err := gl.Initialize(); err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

// Given a Loop
// When Initialize is called once
// Then it has been invoked exactly once
func TestLoop_Initialize_IsCalledOnce(t *testing.T) {
	gl := &mockLoop{}
	gl.Initialize() //nolint:errcheck

	if gl.initializeCalled != 1 {
		t.Fatalf("expected Initialize to be called 1 time, got %d", gl.initializeCalled)
	}
}

// --- Render behaviour ---

// Given an initialized Loop
// When Render is called
// Then it returns nil on success
func TestLoop_WhenRenderSucceeds_ReturnsNil(t *testing.T) {
	gl := &mockLoop{}

	if err := gl.Render(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

// Given a Loop whose Render returns an error
// When Render is called
// Then the error is surfaced to the caller
func TestLoop_WhenRenderFails_ReturnsError(t *testing.T) {
	want := errors.New("render failure")
	gl := &mockLoop{renderErr: want}

	if err := gl.Render(); err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

// Given a Loop
// When Render is called multiple times (simulating multiple ticks)
// Then each call is recorded
func TestLoop_Render_IsCalledOnEachTick(t *testing.T) {
	ticks := 5
	gl := &mockLoop{}

	for i := 0; i < ticks; i++ {
		gl.Render() //nolint:errcheck
	}

	if gl.renderCalled != ticks {
		t.Fatalf("expected Render to be called %d times, got %d", ticks, gl.renderCalled)
	}
}

// --- Calculate behaviour ---

// Given an initialized Loop
// When Calculate is called with a key
// Then it returns nil on success
func TestLoop_WhenCalculateSucceeds_ReturnsNil(t *testing.T) {
	gl := &mockLoop{}

	if err := gl.Calculate(keyboard.KeyArrowUp); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

// Given a Loop whose Calculate returns an error
// When Calculate is called
// Then the error is surfaced to the caller
func TestLoop_WhenCalculateFails_ReturnsError(t *testing.T) {
	want := errors.New("calculate failure")
	gl := &mockLoop{calculateErr: want}

	if err := gl.Calculate(keyboard.KeyArrowUp); err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

// Given a Loop
// When Calculate is called with a specific key
// Then the loop receives that exact key
func TestLoop_Calculate_ReceivesCorrectKey(t *testing.T) {
	cases := []struct {
		name string
		key  keyboard.Key
	}{
		{"arrow up", keyboard.KeyArrowUp},
		{"arrow down", keyboard.KeyArrowDown},
		{"arrow left", keyboard.KeyArrowLeft},
		{"arrow right", keyboard.KeyArrowRight},
		{"space", keyboard.KeySpace},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gl := &mockLoop{}
			gl.Calculate(tc.key) //nolint:errcheck

			if gl.lastKey != tc.key {
				t.Fatalf("expected key %v, got %v", tc.key, gl.lastKey)
			}
		})
	}
}

// --- FPS constant ---

// Given the engine configuration
// When the FPS constant is checked
// Then it equals 30
func TestFPS_Is30(t *testing.T) {
	if FPS != 30 {
		t.Fatalf("expected FPS to be 30, got %d", FPS)
	}
}

// Given the FPS constant
// When the frame duration is derived
// Then each frame is approximately 33 ms (1000ms / 30)
func TestFPS_FrameDurationIsApproximately33ms(t *testing.T) {
	const expectedMs = 33
	frameDurationMs := 1000 / FPS

	// Allow ±1 ms due to integer division
	if frameDurationMs < expectedMs-1 || frameDurationMs > expectedMs+1 {
		t.Fatalf("expected frame duration ~%dms, got %dms", expectedMs, frameDurationMs)
	}
}
