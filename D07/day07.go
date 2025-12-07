package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Point struct{ x, y int }

var directions = []Point{{-1, 0}, {1, 0}}

func readFileRune(fname string) ([][]rune, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var puzzle [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		puzzle = append(puzzle, []rune(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return puzzle, nil
}

func countPaths(grid [][]rune, visited [][]bool, current Point, memo map[Point]int) int {
	if val, ok := memo[current]; ok {
		return val
	}

	if current.y == len(grid)-1 {
		return 1
	}

	visited[current.y][current.x] = true
	count := 0

	if grid[current.y][current.x] == '^' {
		for _, dir := range directions {
			next := Point{current.x + dir.x, current.y + 1}
			if next.x >= 0 && next.x < len(grid[0]) {
				count += countPaths(grid, visited, next, memo)
			}
		}
	} else {
		next := Point{current.x, current.y + 1}
		count += countPaths(grid, visited, next, memo)
	}

	visited[current.y][current.x] = false
	memo[current] = count
	return count
}

func solvePart1(puzzle [][]rune) int {
	result := 0
	for rowIndex := 1; rowIndex < len(puzzle)-1; rowIndex++ {
		for colIndex, colValue := range puzzle[rowIndex] {
			if colValue == '|' {
				if puzzle[rowIndex+1][colIndex] == '^' {
					if colIndex+1 < len(puzzle[0]) {
						puzzle[rowIndex+1][colIndex+1] = '|'
					}
					if colIndex-1 >= 0 {
						puzzle[rowIndex+1][colIndex-1] = '|'
					}
					result++
				} else {
					puzzle[rowIndex+1][colIndex] = '|'
				}
			}
		}
	}
	return result
}

func solvePart2(puzzle [][]rune) int {
	start := Point{}
	for index, r := range puzzle[0] {
		if r == 'S' {
			start = Point{index, 0}
			break
		}
	}

	visited := make([][]bool, len(puzzle))
	for i := range visited {
		visited[i] = make([]bool, len(puzzle[0]))
	}

	memo := make(map[Point]int)
	return countPaths(puzzle, visited, start, memo)
}

func main() {
	timeStart := time.Now()
	// INPUT := "sample.txt"
	INPUT := "input.txt"

	fileContent, err := readFileRune(INPUT)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	sumPart1 := solvePart1(fileContent)
	sumPart2 := solvePart2(fileContent)

	fmt.Println("Part 1:", sumPart1)
	fmt.Println("Part 2:", sumPart2)
	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
