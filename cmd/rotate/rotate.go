package main

import (
	"fmt"
	"os"

	rotate "github.com/Quietus/rotate_exercise"
)

func usage() {
	fmt.Println("Usage: rotate <left|right> <input_file> <output_file>")
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}
	if os.Args[1] != "left" && os.Args[1] != "right" {
		usage()
		os.Exit(1)
	}
	inputFile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Println("Error opening input file:", err)
		os.Exit(1)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = rotate.RotateFile(inputFile, outputFile, os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
