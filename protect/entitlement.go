// Package protect defines the shared EntitlementChecker interface used by
// both the tree-walk interpreter (lang/eval) and the bytecode VM (compiler)
// to enforce #protected[N] tags.
//
// This package has zero dependencies on any other HE package — it's
// deliberately thin so both lang/eval and compiler can import it without
// creating an import cycle.
package protect

import "sync"

// Checker answers "is this #protected[N] tag currently granted?"
// for a given tag name ("protected", "protected1", "protected2", ...).
//
// Implementations must fail closed: any internal error, ambiguity, or
// unknown tag should result in granted=false, not a panic or silent true.
type Checker interface {
	Check(tag string) (granted bool, err error)
}

// StubChecker is an in-memory Checker with no network activity.
// Default behavior: deny everything (fail closed).
// Per-tag overrides and a global default can be set for testing.
type StubChecker struct {
	mu           sync.RWMutex
	defaultGrant bool
	overrides    map[string]bool
}

func NewStubChecker() *StubChecker {
	return &StubChecker{overrides: map[string]bool{}}
}

func (s *StubChecker) SetDefaultGrant(grant bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.defaultGrant = grant
}

func (s *StubChecker) Allow(tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.overrides[tag] = true
}

func (s *StubChecker) Deny(tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.overrides[tag] = false
}

func (s *StubChecker) ClearOverride(tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.overrides, tag)
}

func (s *StubChecker) Check(tag string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.overrides[tag]; ok {
		return v, nil
	}
	return s.defaultGrant, nil
}

// AlwaysGrant grants every tag. For local testing only.
type AlwaysGrant struct{}

func (AlwaysGrant) Check(tag string) (bool, error) { return true, nil }

// AlwaysDeny denies every tag. This is the correct default for any
// unconfigured runtime — fail closed.
type AlwaysDeny struct{}

func (AlwaysDeny) Check(tag string) (bool, error) { return false, nil }
