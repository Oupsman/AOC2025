package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func Solve(part2 bool, inputs []string) int {
	var dir int
	pos := 50
	count := 0
	for _, line := range inputs {
		prev_position := pos
		direction := line[0]
		if direction == 76 {
			dir = -1
		} else if direction == 82 {
			dir = 1
		}
		distance, _ := strconv.Atoi(line[1:])
		if part2 && distance >= 100 {
			count += distance / 100
		}
		pos += (distance % 100) * dir
		if pos > 99 {
			pos -= 100
			if part2 && pos != 0 && prev_position != 0 {
				count++
				fmt.Println("count++")
			}
		}
		if pos < 0 {
			pos += 100
			if part2 && pos != 0 && prev_position != 0 {
				count++
				fmt.Println("count++")
			}

		}
		if pos == 0 {
			count++
		}
		//		fmt.Println("The dial is rotated ", line, "and points at ", pos)
	}
	return count
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
