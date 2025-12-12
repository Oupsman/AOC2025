package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Grid [][]rune

type Present struct {
	shape Grid
	area  int
}

type Region struct {
	x, y    int
	numbers []int
}

func readFile(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var puzzle []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		puzzle = append(puzzle, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return puzzle, nil
}

func solvePart1(input []string) int {
	result := 0
	var presents [6]Present
	var regions []Region
	presentNumber := 0
	for _, line := range input {
		if strings.Contains(line, ":") && !strings.Contains(line, "x") {
			presentNumber, _ = strconv.Atoi(strings.Split(line, ":")[0])
		} else if strings.Contains(line, "#") {
			presents[presentNumber].shape = append(presents[presentNumber].shape, []rune(line))
		} else if strings.Contains(line, "x") {
			fmt.Println("Region detected", line)
			halfs := strings.Split(line, ":")
			dimensions := strings.Split(halfs[0], "x")
			x1, _ := strconv.Atoi(dimensions[0])
			y1, _ := strconv.Atoi(dimensions[1])

			numbersOfPresentsT := strings.Fields(halfs[1])
			numbersOfPresents := make([]int, len(numbersOfPresentsT))
			for i, v := range numbersOfPresentsT {
				numbersOfPresents[i], _ = strconv.Atoi(v)
			}
			regions = append(regions, Region{x: x1, y: y1, numbers: numbersOfPresents})
		}
	}

	for i, present := range presents {
		for _, row := range present.shape {
			for _, col := range row {
				if col == '#' {
					presents[i].area++
				}
			}
		}
	}

	for _, region := range regions {
		totalRegionArea := region.x * region.y
		totalGiftsArea := 0
		for i, v := range region.numbers {
			fmt.Println("Before:", totalGiftsArea)
			totalGiftsArea += presents[i].area * v
			// totalGiftsArea += v * 9
			fmt.Println("After:", totalGiftsArea, i, presents[i].area, v)
		}
		fmt.Println(totalGiftsArea, totalRegionArea)
		if totalGiftsArea < totalRegionArea {
			result++
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
	sumPart1 := solvePart1(fileContent)
	//sumPart2 := solvePart2(fileContent)
	fmt.Println("Part 1:", sumPart1)
	//fmt.Println("Part 2:", sumPart2)

	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
