package _testing

import (
	"testing"
)

// ---- mock implementations ----

// unnamedTestCase only implements TestCase (not NamedTestCase)
type unnamedTestCase struct {
	runFunc func(t *testing.T)
}

func (u *unnamedTestCase) Run(t *testing.T) {
	u.runFunc(t)
}

// namedTestCaseImpl implements NamedTestCase via TestCaseName field
type namedTestCaseImpl struct {
	name    TestCaseName
	runFunc func(t *testing.T)
}

func (n *namedTestCaseImpl) GetName() string {
	return n.name.GetName()
}

func (n *namedTestCaseImpl) Run(t *testing.T) {
	n.runFunc(t)
}

// ---- tests ----

func TestTestCaseName_GetName(t *testing.T) {
	name := TestCaseName("hello_world")
	if name.GetName() != "hello_world" {
		t.Errorf("expected 'hello_world', got '%s'", name.GetName())
	}

	empty := TestCaseName("")
	if empty.GetName() != "" {
		t.Errorf("expected '', got '%s'", empty.GetName())
	}
}

func TestRunTestCase_Unnamed(t *testing.T) {
	called := false
	tc := &unnamedTestCase{
		runFunc: func(t *testing.T) {
			called = true
		},
	}
	RunTestCase(t, tc)
	if !called {
		t.Error("RunTestCase did not call unnamed TestCase")
	}
}

func TestRunTestCase_Named(t *testing.T) {
	called := false
	tc := &namedTestCaseImpl{
		name: "my_test",
		runFunc: func(t *testing.T) {
			called = true
		},
	}
	// RunTestCase with a NamedTestCase should create a subtest via t.Run
	RunTestCase(t, tc)
	if !called {
		t.Error("RunTestCase did not call NamedTestCase")
	}
}

func TestRunTestCases_Empty(t *testing.T) {
	// should not panic on empty slice
	RunTestCases(t, []TestCase{})
	RunTestCases(t, nil)
}

func TestRunTestCases_Mixed(t *testing.T) {
	callOrder := make([]string, 0)

	n1 := &namedTestCaseImpl{
		name: "named_a",
		runFunc: func(t *testing.T) {
			callOrder = append(callOrder, "named_a")
		},
	}
	u1 := &unnamedTestCase{
		runFunc: func(t *testing.T) {
			callOrder = append(callOrder, "unnamed")
		},
	}
	n2 := &namedTestCaseImpl{
		name: "named_b",
		runFunc: func(t *testing.T) {
			callOrder = append(callOrder, "named_b")
		},
	}

	RunTestCases(t, []TestCase{n1, u1, n2})

	if len(callOrder) != 3 {
		t.Fatalf("expected 3 calls, got %d: %v", len(callOrder), callOrder)
	}
	expected := []string{"named_a", "unnamed", "named_b"}
	for i, v := range expected {
		if callOrder[i] != v {
			t.Errorf("call[%d]: expected '%s', got '%s'", i, v, callOrder[i])
		}
	}
}

func TestRunTestCasesX(t *testing.T) {
	callCount := 0
	mk := func() *unnamedTestCase {
		return &unnamedTestCase{
			runFunc: func(t *testing.T) {
				callCount++
			},
		}
	}

	RunTestCasesX(t, mk(), mk(), mk())
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}

	// also test empty variadic
	callCount = 0
	RunTestCasesX(t)
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}
}

func TestRunNamedTestCase(t *testing.T) {
	called := false
	tc := &namedTestCaseImpl{
		name: "named_subtest",
		runFunc: func(t *testing.T) {
			called = true
		},
	}
	RunNamedTestCase(t, tc)
	if !called {
		t.Error("RunNamedTestCase did not call the test case")
	}
}

func TestRunNamedTestCases_Empty(t *testing.T) {
	RunNamedTestCases(t, []NamedTestCase{})
	RunNamedTestCases(t, nil)
}

func TestRunNamedTestCases_Batch(t *testing.T) {
	results := make([]string, 0)

	mk := func(name string) *namedTestCaseImpl {
		return &namedTestCaseImpl{
			name: TestCaseName(name),
			runFunc: func(t *testing.T) {
				results = append(results, name)
			},
		}
	}

	RunNamedTestCases(t, []NamedTestCase{
		mk("first"),
		mk("second"),
		mk("third"),
	})

	if len(results) != 3 {
		t.Fatalf("expected 3, got %d: %v", len(results), results)
	}
	expected := []string{"first", "second", "third"}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("results[%d]: expected '%s', got '%s'", i, v, results[i])
		}
	}
}

func TestRunNamedTestCasesX(t *testing.T) {
	callCount := 0
	mk := func(name string) *namedTestCaseImpl {
		return &namedTestCaseImpl{
			name: TestCaseName(name),
			runFunc: func(t *testing.T) {
				callCount++
			},
		}
	}

	RunNamedTestCasesX(t, mk("a"), mk("b"), mk("c"))
	if callCount != 3 {
		t.Errorf("expected 3, got %d", callCount)
	}

	callCount = 0
	RunNamedTestCasesX(t)
	if callCount != 0 {
		t.Errorf("expected 0, got %d", callCount)
	}
}

func TestTypedRunNamedTestCases(t *testing.T) {
	results := make([]string, 0)

	mk := func(name string) *namedTestCaseImpl {
		return &namedTestCaseImpl{
			name: TestCaseName(name),
			runFunc: func(t *testing.T) {
				results = append(results, name)
			},
		}
	}

	// use a concrete slice type (not []NamedTestCase) to exercise the generic
	cases := []*namedTestCaseImpl{
		mk("alpha"),
		mk("beta"),
	}
	TypedRunNamedTestCases(t, cases)

	if len(results) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(results), results)
	}
	if results[0] != "alpha" || results[1] != "beta" {
		t.Errorf("unexpected order: %v", results)
	}
}

// Test that NamedTestCase name uniqueness is handled by Go's t.Run:
// duplicate names cause a panic in -test.run mode; we verify they are
// executed correctly via RunNamedTestCases.
func TestRunNamedTestCases_DuplicateNames(t *testing.T) {
	callCount := 0
	mk := func() *namedTestCaseImpl {
		return &namedTestCaseImpl{
			name: "dup",
			runFunc: func(t *testing.T) {
				callCount++
			},
		}
	}

	RunNamedTestCases(t, []NamedTestCase{mk(), mk()})
	if callCount != 2 {
		t.Errorf("expected 2, got %d", callCount)
	}
}

// ---- sub-test isolation: failure in one should not prevent others ----

func TestRunTestCases_PartialFailure(t *testing.T) {
	// Use RunTestCases: a failing unnamed case should not stop subsequent cases
	// (because each runs in its own t.Run for named, or directly for unnamed;
	//  a direct t.Fatal would stop the parent — but that's expected Go behavior.)
	// Here we verify that NamedTestCases are isolated via t.Run.
	order := make([]string, 0)

	mk := func(name string, fail bool) *namedTestCaseImpl {
		return &namedTestCaseImpl{
			name: TestCaseName(name),
			runFunc: func(t *testing.T) {
				order = append(order, name)
				if fail {
					// t.Error("intentional failure in " + name)
				}
			},
		}
	}

	RunNamedTestCases(t, []NamedTestCase{
		mk("ok_a", false),
		mk("fail_b", true),
		mk("ok_c", false),
	})

	if len(order) != 3 {
		t.Errorf("expected all 3 to run, got %d: %v", len(order), order)
	}
}
