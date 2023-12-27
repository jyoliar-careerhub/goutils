package ptr

import (
	"errors"
	"testing"
)

type SampleError struct{}

func (e SampleError) Error() string {
	return "sample error struct"
}

func TestIsNilErr(t *testing.T) {
	t.Run("basic error", func(t *testing.T) {
		type TestCase struct {
			testName string
			target   error
			expected bool
		}

		var sampleErr *SampleError = nil
		testCases := []TestCase{
			{testName: "nil", target: nil, expected: true},
			{testName: "errors.New(\"sample error\")", target: errors.New("sample error"), expected: false},
			{testName: "&SampleError{}", target: &SampleError{}, expected: false},
			{testName: "var sampleErr *SampleError = nil", target: sampleErr, expected: true},
		}

		for _, testCase := range testCases {
			t.Run(testCase.testName, func(t *testing.T) {
				isNil := IsNil(testCase.target)
				if isNil != testCase.expected {
					t.Errorf("expected: %v, actual: %v", testCase.expected, isNil)
				}
			})
		}
	})

	t.Run("custom pointer error", func(t *testing.T) {
		type TestCase struct {
			testName string
			target   *SampleError
			expected bool
		}

		var sampleErr *SampleError = nil
		testCases := []TestCase{
			{testName: "nil", target: nil, expected: true},
			{testName: "&SampleError{}", target: &SampleError{}, expected: false},
			{testName: "var sampleErr *SampleError = nil", target: sampleErr, expected: true},
		}

		for _, testCase := range testCases {
			t.Run(testCase.testName, func(t *testing.T) {
				actual := IsNil(testCase.target)
				if actual != testCase.expected {
					t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
				}
			})
		}
	})

	t.Run("custom strict error", func(t *testing.T) {
		type TestCase struct {
			testName string
			target   SampleError
			expected bool
		}

		testCases := []TestCase{
			{testName: "SampleError{}", target: SampleError{}, expected: false},
			{testName: "Nothing", expected: false},
		}

		for _, testCase := range testCases {
			t.Run(testCase.testName, func(t *testing.T) {
				actual := IsNil(testCase.target)
				if actual != testCase.expected {
					t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
				}
			})
		}
	})
}
