// Package eval — entitlement checking for #protected[N] tags.
//
// This file re-exports the protect.Checker interface and its stub
// implementations under the eval package's own names, preserving
// backward compatibility with existing code that references
// eval.EntitlementChecker / eval.StubChecker / etc.
//
// The canonical definitions now live in hunterlang/protect.
package eval

import "hunterlang/protect"

// EntitlementChecker is an alias for protect.Checker.
// Kept here so nothing in lang/eval needs to change its import paths.
type EntitlementChecker = protect.Checker

// StubChecker re-exports protect.StubChecker.
type StubChecker = protect.StubChecker

// AlwaysGrant re-exports protect.AlwaysGrant.
type AlwaysGrant = protect.AlwaysGrant

// AlwaysDeny re-exports protect.AlwaysDeny.
type AlwaysDeny = protect.AlwaysDeny

// NewStubChecker re-exports protect.NewStubChecker.
func NewStubChecker() *protect.StubChecker {
	return protect.NewStubChecker()
}
