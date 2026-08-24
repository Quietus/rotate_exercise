package rotate

import (
	"testing"
)

func TestLeftRotate(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		rotations int
		expected  []byte
	}{
		{"Case 1", []byte{2, 3, 4, 5}, 1, []byte{4, 6, 8, 10}},
		{"Case 2", []byte{255, 0, 0, 0}, 1, []byte{254, 0, 0, 1}},
		{"Case 3", []byte{128, 128, 128, 128}, 1, []byte{1, 1, 1, 1}},
	}
	for _, tc := range testCases {
		result := LeftRotate(tc.input)
		if result[0] != tc.expected[0] || result[1] != tc.expected[1] || result[2] != tc.expected[2] || result[3] != tc.expected[3] {
			t.Errorf("LeftRotate failed for %s: got %v: expected: %v", tc.name, result, tc.expected)
		}
	}
}

func TestRightRotate(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		rotations int
		expected  []byte
	}{
		{"Case 1", []byte{2, 3, 4, 5}, 1, []byte{129, 1, 130, 2}},
		{"Case 2", []byte{255, 0, 0, 0}, 1, []byte{127, 128, 0, 0}},
		{"Case 3", []byte{128, 128, 128, 128}, 1, []byte{64, 64, 64, 64}},
	}
	for _, tc := range testCases {
		result := RightRotate(tc.input)
		if result[0] != tc.expected[0] || result[1] != tc.expected[1] || result[2] != tc.expected[2] || result[3] != tc.expected[3] {
			t.Errorf("RightRotate failed for %s: got %v: expected: %v", tc.name, result, tc.expected)
		}
	}
}
