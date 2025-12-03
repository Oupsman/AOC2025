package main

import (
	"bufio"
	"fmt"
	"os"
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

func findHighestCombinaision(bank string, number, initialpos int) int {
	n := len(bank)
	if number <= 0 || initialpos+number > n {
		return 0
	}

	maxPos := initialpos
	maxVal := bank[initialpos]
	for j := initialpos; j <= n-number; j++ {
		if bank[j] > maxVal {
			maxVal = bank[j]
			maxPos = j
		}
	}

	return int(maxVal-'0')*intPow(10, number-1) + findHighestCombinaision(bank, number-1, maxPos+1)
}

func intPow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func Solve(part2 bool, inputs []string) int {
	var result int = 0
	for _, bank := range inputs {
		voltage := 0
		if !part2 {
			voltage = findHighestCombinaision(bank, 2, 0)
		} else {
			voltage = findHighestCombinaision(bank, 12, 0)
		}
		result += voltage
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
