package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func checkID(ingredientID, range1, range2 int) bool {
	return ingredientID >= range1 && ingredientID <= range2
}

func Solve(part2 bool, inputs []string) int64 {
	var result int64 = 0
	var ranges [][2]int
	var ingredients []int
	var lineCounter int = 0
	var line string
	line = inputs[lineCounter]
	for line != "" {

		rangeValue := strings.Split(line, "-")
		value1, _ := strconv.Atoi(rangeValue[0])
		value2, _ := strconv.Atoi(rangeValue[1])
		ranges = append(ranges, [2]int{value1, value2})

		lineCounter++
		line = inputs[lineCounter]

	}

	for line == "" {
		lineCounter++
		line = inputs[lineCounter]
	}

	for line != "" && lineCounter < len(inputs) {
		line = inputs[lineCounter]

		value, _ := strconv.Atoi(line)
		ingredients = append(ingredients, value)
		lineCounter++

	}
	if part2 {
		var mergedRanges [][2]int
		// Merge overlapping ranges
		fmt.Println("Ranges:", ranges)
		// Sort the ranges by start value
		sort.Slice(ranges, func(i, j int) bool {
			return ranges[i][0] < ranges[j][0]
		})
		mergedRanges = append(mergedRanges, ranges[0])
		for _, interval := range ranges {
			last := &mergedRanges[len(mergedRanges)-1]
			if interval[0] <= last[1] {
				if interval[1] > last[1] {
					last[1] = interval[1]
				}
			} else {
				mergedRanges = append(mergedRanges, interval)
			}
		}
		fmt.Println("Merged ranges:", mergedRanges)
		for _, rangeValues := range mergedRanges {
			result += int64(rangeValues[1] - rangeValues[0] + 1)
		}

		return result
	}
	for _, ingredientID := range ingredients {
		for _, rangeValue := range ranges {
			if checkID(ingredientID, rangeValue[0], rangeValue[1]) {
				result++
				break
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
