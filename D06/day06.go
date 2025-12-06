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
func reconstructNumbersFromColumn(numbers ...int) []int {

	strNumbers := make([]string, len(numbers))
	for i, num := range numbers {
		strNumbers[i] = strconv.Itoa(num)
	}

	maxLen := 0
	for _, s := range strNumbers {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}

	padded := make([]string, len(strNumbers))
	for i, s := range strNumbers {
		padded[i] = s + strings.Repeat(" ", maxLen-len(s))
	}
	var reconstructed []int
	for col := maxLen - 1; col >= 0; col-- {
		var numStr string
		for _, s := range padded {
			if col < len(s) && s[col] != ' ' {
				numStr += string(s[col])
			}
		}
		if numStr != "" {
			num, _ := strconv.Atoi(numStr)
			reconstructed = append(reconstructed, num)
		}
	}

	for i, j := 0, len(reconstructed)-1; i < j; i, j = i+1, j-1 {
		reconstructed[i], reconstructed[j] = reconstructed[j], reconstructed[i]
	}

	return reconstructed
}

func Solve(part2 bool, inputs []string) int64 {
	var result int64 = 0
	var operators []string
	var numbers [][]int
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
	// fmt.Println(numbers)
	if !part2 {
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
		for i := len(numbers[0]) - 1; i >= 0; i-- {
			var colNumbers []int
			for _, col := range numbers {
				colNumbers = append(colNumbers, col[i])
			}
			fmt.Println("Origin numbers", colNumbers)
			numbersPart2 := reconstructNumbersFromColumn(colNumbers...)
			fmt.Println("Interpreted:", numbersPart2)
			if len(numbersPart2) == 0 {
				continue
			}

			op := operators[i]

			resultLine := numbersPart2[0]
			for i := 1; i < len(numbersPart2); i++ {
				switch op {
				case "*":
					fmt.Println("Multiplying")
					resultLine *= numbersPart2[i]
				case "+":
					fmt.Println("Adding")
					resultLine += numbersPart2[i]
				}
			}
			fmt.Println("Result of line: ", resultLine)
			result += int64(resultLine)
		}

	}
	return result
}

func main() {
	timeStart := time.Now()

	INPUT := "sample.txt"
	// INPUT := "input.txt"
	fileContent := readFile(INPUT)
	sumPart1 := Solve(false, fileContent)
	sumPart2 := Solve(true, fileContent)

	fmt.Println("Part1:", sumPart1)
	fmt.Println("Part2:", sumPart2)

	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
