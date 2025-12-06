package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func readFile(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return lines, nil
}

func reconstructNumbers(input []string) []int {
	maxCol := 0
	for _, line := range input {
		if len(line) > maxCol {
			maxCol = len(line)
		}
	}

	padded := make([]string, len(input))
	for i, line := range input {
		padded[i] = line + strings.Repeat(" ", maxCol-len(line))
	}

	var numbers []int
	for j := 0; j < maxCol; j++ {
		var numStr strings.Builder
		for i := 0; i < len(padded); i++ {
			if j < len(padded[i]) && padded[i][j] != ' ' {
				numStr.WriteByte(padded[i][j])
			}
		}
		if numStr.Len() > 0 {
			number, _ := strconv.Atoi(numStr.String())
			numbers = append(numbers, number)
		}
	}
	return numbers
}

func doOperation(numbers []int, operator rune) int64 {
	if len(numbers) == 0 {
		return 0
	}
	partialResult := int64(numbers[0])
	for i := 1; i < len(numbers); i++ {
		switch operator {
		case '*':
			partialResult *= int64(numbers[i])
		case '+':
			partialResult += int64(numbers[i])
		}
	}
	return partialResult
}

func Solve(part2 bool, inputs []string) int64 {
	if len(inputs) < 2 {
		return 0
	}

	var result int64
	operatorsLine := inputs[len(inputs)-1]

	if !part2 {
		reOperators := regexp.MustCompile(`[+*]`)
		operators := reOperators.FindAllString(operatorsLine, -1)
		if len(operators) == 0 {
			return 0
		}

		var numbers [][]int
		for i := 0; i < len(inputs)-1; i++ {
			re := regexp.MustCompile(`\d+`)
			var numbersLine []int
			for _, value := range re.FindAllString(inputs[i], -1) {
				number, _ := strconv.Atoi(value)
				numbersLine = append(numbersLine, number)
			}
			numbers = append(numbers, numbersLine)
		}

		for col := 0; col < len(numbers[0]) && col < len(operators); col++ {
			op := operators[col][0]
			var colNumbers []int
			for row := 0; row < len(numbers); row++ {
				colNumbers = append(colNumbers, numbers[row][col])
			}
			result += doOperation(colNumbers, rune(op))
		}
	} else {
		var opIndices []int
		for i, char := range operatorsLine {
			if char == '*' || char == '+' {
				opIndices = append(opIndices, i)
			}
		}

		for idx, start := range opIndices {
			end := len(inputs[0])
			if idx < len(opIndices)-1 {
				end = opIndices[idx+1]
			}

			var numStrings []string
			for i := 0; i < len(inputs)-1; i++ {
				line := inputs[i]
				if start >= len(line) {
					continue
				}
				if end > len(line) {
					end = len(line)
				}
				numStrings = append(numStrings, line[start:end])
			}

			numbersFromString := reconstructNumbers(numStrings)
			op := rune(operatorsLine[start])
			result += doOperation(numbersFromString, op)
		}
	}

	return result
}

func main() {
	timeStart := time.Now()
	// INPUT := "sample.txt"
	INPUT := "input.txt"

	fileContent, err := readFile(INPUT)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	sumPart1 := Solve(false, fileContent)
	sumPart2 := Solve(true, fileContent)

	fmt.Println("Part 1:", sumPart1)
	fmt.Println("Part 2:", sumPart2)
	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
