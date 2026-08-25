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
	defer inputFile.Close()

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
	defer leftOutputFile.Close()

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
	defer rightOutputFile.Close()

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
