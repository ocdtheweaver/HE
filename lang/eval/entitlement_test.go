package eval

import (
	"testing"

	"hunterlang/lang/lexer"
	"hunterlang/lang/parser"
)

// These tests exercise the Stage A / Step 2 stub in isolation, before
// any enforcement code (Step 3) exists to call it. The point is to lock
// in the fail-closed contract now, so enforcement can be built against
// a checker whose behavior is already proven correct.

func TestAlwaysDeny(t *testing.T) {
	var c EntitlementChecker = AlwaysDeny{}
	granted, err := c.Check("protected")
	if err != nil {
		t.Fatalf("AlwaysDeny.Check returned an error: %v", err)
	}
	if granted {
		t.Fatal("AlwaysDeny granted a tag — it must never grant anything")
	}
}

func TestAlwaysGrant(t *testing.T) {
	var c EntitlementChecker = AlwaysGrant{}
	granted, err := c.Check("protected1")
	if err != nil {
		t.Fatalf("AlwaysGrant.Check returned an error: %v", err)
	}
	if !granted {
		t.Fatal("AlwaysGrant denied a tag — it must never deny anything")
	}
}

func TestStubChecker_DefaultsToDeny(t *testing.T) {
	s := NewStubChecker()
	granted, err := s.Check("protected")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if granted {
		t.Fatal("a freshly-constructed StubChecker must default to deny (fail closed)")
	}
}

func TestStubChecker_SetDefaultGrant(t *testing.T) {
	s := NewStubChecker()
	s.SetDefaultGrant(true)
	granted, _ := s.Check("anything")
	if !granted {
		t.Fatal("SetDefaultGrant(true) should cause un-overridden tags to be granted")
	}

	s.SetDefaultGrant(false)
	granted, _ = s.Check("anything")
	if granted {
		t.Fatal("SetDefaultGrant(false) should cause un-overridden tags to be denied")
	}
}

func TestStubChecker_PerTagOverrides(t *testing.T) {
	s := NewStubChecker() // default: deny
	s.Allow("protected1")
	s.Deny("protected2")

	cases := []struct {
		tag     string
		want    bool
		comment string
	}{
		{"protected", false, "un-overridden tag should fall back to default (deny)"},
		{"protected1", true, "explicitly allowed tag should be granted"},
		{"protected2", false, "explicitly denied tag should stay denied"},
	}

	for _, c := range cases {
		got, err := s.Check(c.tag)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", c.tag, err)
		}
		if got != c.want {
			t.Fatalf("%s: got %v, want %v — %s", c.tag, got, c.want, c.comment)
		}
	}
}

func TestStubChecker_ClearOverride(t *testing.T) {
	s := NewStubChecker()
	s.SetDefaultGrant(false)
	s.Allow("protected1")

	granted, _ := s.Check("protected1")
	if !granted {
		t.Fatal("expected protected1 to be granted via override before clearing")
	}

	s.ClearOverride("protected1")
	granted, _ = s.Check("protected1")
	if granted {
		t.Fatal("after ClearOverride, protected1 should fall back to the default (deny)")
	}
}

func TestRuntime_DefaultsToAlwaysDeny(t *testing.T) {
	r := newRuntime()
	granted, err := r.entitlement.Check("protected")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if granted {
		t.Fatal("a newly constructed Runtime must default to AlwaysDeny — fail closed")
	}
}

func TestInterpreter_SetEntitlementChecker(t *testing.T) {
	i := NewInterpreter()
	stub := NewStubChecker()
	stub.SetDefaultGrant(true)
	i.SetEntitlementChecker(stub)

	if err := i.Run(nil); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	granted, err := i.runtime.entitlement.Check("protected")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !granted {
		t.Fatal("Interpreter.SetEntitlementChecker should have propagated into the runtime")
	}
}

// ── End-to-end enforcement tests (Stage A, Step 3) ──────────────────────────
//
// The tests above exercise the checker in isolation. These run real HE
// source through a full Interpreter with a configured checker, proving
// that callMethod actually consults it — not just that the interface
// exists.

func runHE(t *testing.T, src string, checker EntitlementChecker) error {
	t.Helper()
	lx := lexer.New(src)
	p := parser.New(lx)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	i := NewInterpreter()
	if checker != nil {
		i.SetEntitlementChecker(checker)
	}
	return i.Run(prog)
}

func TestEnforcement_DeniedByDefault(t *testing.T) {
	src := `
create Vault #protected1 [
  can: [
    open [
      say "opened"
    ]
  ]
]
tell Vault to open
`
	err := runHE(t, src, nil) // nil → defaults to AlwaysDeny
	if err == nil {
		t.Fatal("expected an error from a denied #protected1 call, got nil")
	}
}

func TestEnforcement_GrantedTagAllowsCall(t *testing.T) {
	stub := NewStubChecker()
	stub.Allow("protected1")

	src := `
create Vault #protected1 [
  can: [
    open [
      say "opened"
    ]
  ]
]
tell Vault to open
`
	if err := runHE(t, src, stub); err != nil {
		t.Fatalf("expected granted tag to allow the call, got error: %v", err)
	}
}

func TestEnforcement_PerAbilityTag(t *testing.T) {
	stub := NewStubChecker()
	stub.Deny("protected2")

	src := `
create Store [
  can: [
    unlockSkin #protected2 [
      say "unlocked"
    ]
    browse [
      say "browsing"
    ]
  ]
]
tell Store to browse
tell Store to unlockSkin
`
	err := runHE(t, src, stub)
	if err == nil {
		t.Fatal("expected the #protected2 ability call to be denied")
	}
}

func TestEnforcement_UnprotectedAbilityUnaffected(t *testing.T) {
	// No checker configured at all (defaults to AlwaysDeny) — but since
	// nothing here is tagged, nothing should ever consult the checker.
	src := `
create Catalog [
  can: [
    browse [
      say "browsing"
    ]
  ]
]
tell Catalog to browse
`
	if err := runHE(t, src, nil); err != nil {
		t.Fatalf("unprotected ability should never be denied, got: %v", err)
	}
}

func TestEnforcement_CatchableViaTryOr(t *testing.T) {
	src := `
create Vault #protected1 [
  can: [
    open [
      say "opened"
    ]
  ]
]
set caught to false
try [
  tell Vault to open
] or (err) [
  set caught to true
]
if caught is not true then [
  say "DID NOT CATCH"
]
`
	// nil checker → AlwaysDeny, so the try/or should catch the denial
	// and the program should complete without propagating an error.
	if err := runHE(t, src, nil); err != nil {
		t.Fatalf("denial should have been caught by try/or, but Run still errored: %v", err)
	}
}

func TestEnforcement_WholeObjectTagAppliesToAllAbilities(t *testing.T) {
	stub := NewStubChecker() // default deny
	stub.Allow("protected")

	src := `
make Premium #protected [
  can: [
    featureA [ say "a" ]
    featureB [ say "b" ]
  ]
]
tell Premium to featureA
tell Premium to featureB
`
	if err := runHE(t, src, stub); err != nil {
		t.Fatalf("granting the object-level tag should allow every ability on it: %v", err)
	}
}
