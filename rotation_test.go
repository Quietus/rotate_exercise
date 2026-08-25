package rotate

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"testing"
)

func compareFiles(file1, file2 io.ReadSeeker) (bool, error) {
	_, err := file1.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	_, err = file2.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}

	buf1 := make([]byte, 1024)
	buf2 := make([]byte, 1024)

	for {
		n1, err1 := file1.Read(buf1)
		n2, err2 := file2.Read(buf2)

		if n1 != n2 || (err1 != nil && err1 != io.EOF) || (err2 != nil && err2 != io.EOF) {
			return false, nil
		}

		if n1 == 0 {
			break
		}

		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
	}

	return true, nil
}

func TestRotateFile(t *testing.T) {
	var limit int64 = 1000
	// Create a temporary input file with random data
	inputFile, err := os.CreateTemp("", "input")
	if err != nil {
		t.Fatalf("Failed to create temporary input file: %v", err)
	}
	t.Cleanup(func() {
		inputFile.Close()
		os.Remove(inputFile.Name())
	})

	_, err = io.Copy(inputFile, io.LimitReader(rand.Reader, limit))
	if err != nil {
		t.Fatalf("Failed to write random data to input file: %v", err)
	}
	inputFile.Seek(0, io.SeekStart)

	// Create a temporary output file
	leftOutputFile, err := os.CreateTemp("", "leftoutput")
	if err != nil {
		t.Fatalf("Failed to create temporary output file: %v", err)
	}
	t.Cleanup(func() {
		leftOutputFile.Close()
		os.Remove(leftOutputFile.Name())
	})

	// Test left rotation
	err = RotateFile(inputFile, leftOutputFile, Left)
	if err != nil {
		t.Errorf("Left rotation failed: %v", err)
	}

	equal, err := compareFiles(inputFile, leftOutputFile)
	if err != nil {
		t.Errorf("Failed to compare files: %v", err)
	}
	if equal {
		t.Errorf("Files are equal following left rotation, expected them to differ")
	}

	// Reset the left Output file's read position to the beginning for the next test
	leftOutputFile.Seek(0, io.SeekStart)

	rightOutputFile, err := os.CreateTemp("", "rightoutput")
	if err != nil {
		t.Fatalf("Failed to create temporary output file: %v", err)
	}
	t.Cleanup(func() {
		rightOutputFile.Close()
		os.Remove(rightOutputFile.Name())
	})

	// Test right rotation
	err = RotateFile(leftOutputFile, rightOutputFile, Right)
	if err != nil {
		t.Errorf("Right rotation failed: %v", err)
	}

	equal, err = compareFiles(leftOutputFile, rightOutputFile)
	if err != nil {
		t.Errorf("Failed to compare files: %v", err)
	}
	if equal {
		t.Errorf("Files are equal following right rotation, expected them to differ")
	}

	equal, err = compareFiles(inputFile, rightOutputFile)
	if err != nil {
		t.Errorf("Failed to compare files: %v", err)
	}
	if !equal {
		t.Errorf("Files are not equal after left and right rotations, expected them to be equal")
	}
}

func TestRotateFileInvalidDirection(t *testing.T) {
	input := bytes.NewReader([]byte{0x01, 0x02, 0x03})
	output := bytes.NewBuffer(nil)

	err := RotateFile(input, output, "invalid_direction")
	if err == nil {
		t.Errorf("Expected error for invalid rotation direction, got nil")
	}
}

func TestRotateFileEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		inputData []byte
		expected  []byte
		direction string
	}{
		{
			name:      "Empty file",
			inputData: []byte{},
			expected:  []byte{},
			direction: Left,
		},
		{
			name:      "Single byte file",
			inputData: []byte{0xFF},
			expected:  []byte{0xFF},
			direction: Right,
		},
		{
			name:      "Two byte file left rotation",
			inputData: []byte{0x01, 0x02},
			expected:  []byte{0x02, 0x04},
			direction: Left,
		},
		{
			name:      "Two byte file right rotation",
			inputData: []byte{0x01, 0x02},
			expected:  []byte{0x00, 0x81},
			direction: Right,
		},
		{
			name:      "All zero byte file left rotation",
			inputData: []byte{0x00, 0x00, 0x00},
			expected:  []byte{0x00, 0x00, 0x00},
			direction: Left,
		},
		{
			name:      "All zero byte file right rotation",
			inputData: []byte{0x00, 0x00, 0x00},
			expected:  []byte{0x00, 0x00, 0x00},
			direction: Right,
		},
		{
			name:      "All one byte file left rotation",
			inputData: []byte{0xFF, 0xFF, 0xFF},
			expected:  []byte{0xFF, 0xFF, 0xFF},
			direction: Left,
		},
		{
			name:      "All one byte file right rotation",
			inputData: []byte{0xFF, 0xFF, 0xFF},
			expected:  []byte{0xFF, 0xFF, 0xFF},
			direction: Right,
		},
		{
			name:      "Alternating pattern left rotation",
			inputData: []byte{0xAA, 0x55, 0xAA, 0x55},
			expected:  []byte{0x54, 0xAB, 0x54, 0xAB},
			direction: Left,
		},
		{
			name:      "Alternating pattern right rotation",
			inputData: []byte{0x54, 0xAB, 0x54, 0xAB},
			expected:  []byte{0xAA, 0x55, 0xAA, 0x55},
			direction: Right,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := bytes.NewReader(tc.inputData)
			output := bytes.NewBuffer(nil)

			err := RotateFile(input, output, tc.direction)
			if err != nil {
				t.Fatalf("RotateFile failed: %v", err)
			}

			if !bytes.Equal(output.Bytes(), tc.expected) {
				t.Errorf("Expected output %v, got %v", tc.expected, output.Bytes())
			}
		})
	}
}
