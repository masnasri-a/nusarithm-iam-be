package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Path to the docs.go file
	filePath := "docs/docs.go"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: %s does not exist\n", filePath)
		os.Exit(1)
	}

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip lines containing LeftDelim or RightDelim
		if strings.Contains(line, "LeftDelim:") || strings.Contains(line, "RightDelim:") {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Write the cleaned content back to the file
	output, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	fmt.Println("Successfully cleaned docs.go - removed LeftDelim and RightDelim lines")
}
