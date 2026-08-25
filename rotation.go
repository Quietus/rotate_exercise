package rotate

import (
	"bufio"
	"fmt"
	"io"
)

const (
	// Left rotation direction
	Left = "left"
	// Right rotation direction
	Right = "right"
)

type RotateWriter interface {
	Write(data []byte) error
	Flush() error
}

type LeftRotator struct {
	buffer     byte
	output     io.Writer
	lastBit    int
	lastBitSet bool
}

type RightRotator struct {
	buffer      byte
	output      io.Writer
	firstBit    int
	firstBitSet bool
}

func NewLeftRotator(output io.Writer) *LeftRotator {
	return &LeftRotator{
		buffer:     0,
		output:     output,
		lastBit:    0,
		lastBitSet: false,
	}
}

func (r *LeftRotator) Write(data []byte) error {
	index := 0
	if !r.lastBitSet {
		if data[0] >= 128 {
			r.lastBit = 1
		}
		r.lastBitSet = true
		r.buffer = (data[0] - 128) << 1
		index = 1
	}
	for ; index < len(data); index++ {
		if data[index] >= 128 {
			_, err := r.output.Write([]byte{r.buffer + 1})
			if err != nil {
				return err
			}
			r.buffer = (data[index] - 128) << 1
		} else {
			_, err := r.output.Write([]byte{r.buffer})
			if err != nil {
				return err
			}
			r.buffer = data[index] << 1
		}
	}
	return nil
}

func (r *LeftRotator) Flush() error {
	if r.lastBitSet && r.lastBit == 1 {
		_, err := r.output.Write([]byte{r.buffer + 1})
		if err != nil {
			return err
		}
	} else if r.lastBitSet {
		_, err := r.output.Write([]byte{r.buffer})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRightRotator(output io.Writer) *RightRotator {
	return &RightRotator{
		buffer:      0,
		output:      output,
		firstBit:    0,
		firstBitSet: false,
	}
}

func (r *RightRotator) Write(data []byte) error {
	var carryBit int = 0
	for i := 0; i < len(data); i++ {
		if r.firstBitSet {
			carryBit = r.firstBit
			r.firstBitSet = false
		} else {
			carryBit = int(r.buffer & 1)
		}
		r.buffer = data[i]
		r.output.Write([]byte{(r.buffer >> 1) + byte(carryBit<<7)})
	}
	return nil
}

func (r *RightRotator) Flush() error {
	return nil
}

func RotateFile(inputFile io.ReadSeeker, outputFile io.Writer, direction string) error {
	var firstBit int = 0
	var firstBitSet = false

	if direction != Left && direction != Right {
		return fmt.Errorf("invalid rotation direction: %s", direction)
	}
	if direction == Right {
		_, err := inputFile.Seek(-1, io.SeekEnd)
		if err != nil {
			return fmt.Errorf("failed to reach end of input file: %v", err)
		}
		lastByte := make([]byte, 1)
		_, err = inputFile.Read(lastByte)
		if err != nil {
			return fmt.Errorf("failed to read last byte of input file: %v", err)
		}
		firstBit = int(lastByte[0] & 1)
		_, err = inputFile.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek to start of input file: %v", err)
		}
		firstBitSet = true
	}
	scanner := bufio.NewScanner(bufio.NewReader(inputFile))
	scanner.Split(bufio.ScanBytes)
	bufferedWriter := bufio.NewWriter(outputFile)
	var rotateWriter RotateWriter
	if direction == Left {
		rotateWriter = NewLeftRotator(bufferedWriter)
	} else {
		rightWriter := NewRightRotator(bufferedWriter)
		rightWriter.firstBit = firstBit
		rightWriter.firstBitSet = firstBitSet
		rotateWriter = rightWriter
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		err := rotateWriter.Write(line)
		if err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}
	if err := rotateWriter.Flush(); err != nil {
		return fmt.Errorf("error flushing output file: %v", err)
	}
	bufferedWriter.Flush()
	return nil
}
