package _testing

import "testing"

// todo move to goutils
type TestCaseName string

func (name TestCaseName) GetName() string {
	return string(name)
}

type TestCase interface {
	Run(t *testing.T)
}

type NamedTestCase interface {
	TestCase
	GetName() string
}

func RunTestCase(t *testing.T, tc TestCase) {
	ntc, ok := tc.(NamedTestCase)
	if ok {
		t.Run(ntc.GetName(), tc.Run)
	} else {
		tc.Run(t)
	}
}

func RunTestCases(t *testing.T, tcs []TestCase) {
	for _, tc := range tcs {
		RunTestCase(t, tc)
	}
}

func RunTestCasesX(t *testing.T, tcs ...TestCase) {
	RunTestCases(t, tcs)
}

func RunNamedTestCase(t *testing.T, tc NamedTestCase) {
	t.Run(tc.GetName(), tc.Run)
}

func RunNamedTestCases(t *testing.T, tcs []NamedTestCase) {
	for _, tc := range tcs {
		RunNamedTestCase(t, tc)
	}
}

func RunNamedTestCasesX(t *testing.T, tcs ...NamedTestCase) {
	RunNamedTestCases(t, tcs)
}

func TypedRunNamedTestCases[Cases ~[]Case, Case NamedTestCase](t *testing.T, tcs Cases) {
	for _, tc := range tcs {
		RunNamedTestCase(t, tc)
	}
}
