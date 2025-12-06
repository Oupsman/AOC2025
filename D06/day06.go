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

func readFile(fname string) []string {
	var lines []string
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}

func reconstructNumbers(input []string) []int {
	var numbers []int
	maxCol := 0

	for _, line := range input {
		if len(line) > maxCol {
			maxCol = len(line)
		}
	}

	padded := make([]string, len(input))
	for i, line := range input {
		if len(line) < maxCol {
			padded[i] = line + strings.Repeat(" ", maxCol-len(line))
		} else {
			padded[i] = line
		}
	}

	for i, line := range padded {
		fmt.Printf("Line %d: |%s|\n", i, line)
	}

	for j := 0; j < maxCol; j++ {
		var numString string

		for i := 0; i < len(padded); i++ {
			if j < len(padded[i]) && padded[i][j] != ' ' {
				numString += string(padded[i][j])
			}
		}

		if numString != "" {
			number, _ := strconv.Atoi(numString)
			numbers = append(numbers, number)
		}
	}

	return numbers
}

func doOperation(numbersFromString []int, operator uint8) int64 {
	fmt.Println("Numbers from string: ", numbersFromString)
	partialResult := int64(numbersFromString[0])
	for i := 1; i < len(numbersFromString); i++ {
		number := numbersFromString[i]
		if operator == '*' {
			partialResult *= int64(number)
		} else {
			partialResult += int64(number)
		}
		fmt.Println("Partial result: ", partialResult)
	}
	return partialResult
}

func Solve(part2 bool, inputs []string) int64 {
	var result int64 = 0
	var operators []string
	var numbers [][]int

	// fmt.Println(numbers)
	if !part2 {
		reOperators := regexp.MustCompile(`[\+\*]`)

		for _, char := range reOperators.FindAllString(inputs[(len(inputs)-1)], -1) {
			operators = append(operators, char)
		}
		// fmt.Println(operators)

		for i := 0; i < len(inputs)-1; i++ {
			var numbersLine []int
			re := regexp.MustCompile(`\d+`)
			// result := re.ReplaceAllString(input, " ")
			for _, value := range re.FindAllString(inputs[i], -1) {
				number, _ := strconv.Atoi(value)
				numbersLine = append(numbersLine, number)
			}
			numbers = append(numbers, numbersLine)
		}
		for i := 0; i < len(numbers[0]); i++ {
			var resultLine int64 = 0
			for j := 0; j < len(numbers); j++ {
				// fmt.Println(numbers[j][i])
				if operators[i] == "+" {
					// fmt.Println("Adding")
					resultLine += int64(numbers[j][i])
				} else if operators[i] == "*" {
					// fmt.Println("Multiplying")
					if j != 0 {
						resultLine *= int64(numbers[j][i])
					} else {
						resultLine = int64(numbers[j][i])
					}
				}
				// fmt.Println("ResultLine: ", resultLine)
			}
			// fmt.Println("Result of line: ", resultLine)
			result += resultLine
		}

		return result
	} else {
		var startingColumns []int
		for i, char := range inputs[(len(inputs) - 1)] {
			if char == '*' || char == '+' {
				startingColumns = append(startingColumns, i)
			}
		}
		// Now we cut the other part of the array
		for startingIndex := 0; startingIndex < len(startingColumns); startingIndex++ {
			var numStrings []string
			if startingIndex != len(startingColumns)-1 {
				for i := 0; i < len(inputs)-1; i++ {
					numString := inputs[i][startingColumns[startingIndex] : startingColumns[startingIndex+1]-1]
					fmt.Printf("NumString: |%s|\n", numString)
					numStrings = append(numStrings, numString)
				}
				fmt.Println(numStrings)
				numbersFromString := reconstructNumbers(numStrings)
				operator := inputs[len(inputs)-1][startingColumns[startingIndex]]

				result += doOperation(numbersFromString, operator)
			} else {
				for i := 0; i < len(inputs)-1; i++ {
					numString := inputs[i][startingColumns[startingIndex]:len(inputs[i])]
					fmt.Printf("NumString: |%s|\n", numString)
					numStrings = append(numStrings, numString)
				}
				fmt.Println(numStrings)
				numbersFromString := reconstructNumbers(numStrings)
				operator := inputs[len(inputs)-1][startingColumns[startingIndex]]

				result += doOperation(numbersFromString, operator)
			}
		}
	}
	return result
}

func main() {
	timeStart := time.Now()

	// INPUT := "sample.txt"
	INPUT := "input.txt"
	fileContent := readFile(INPUT)
	sumPart1 := Solve(false, fileContent)
	sumPart2 := Solve(true, fileContent)

	fmt.Println("Part1:", sumPart1)
	fmt.Println("Part2:", sumPart2)

	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
