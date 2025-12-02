package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
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
		lines = strings.Split(line, ",")
	}
	return lines
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func generateInvalidIDs(L, d int) []int64 {
	var ids []int64
	repeatCount := L / d
	start := int64(math.Pow(10, float64(d-1)))
	end := int64(math.Pow(10, float64(d)))

	for prefix := start; prefix < end; prefix++ {
		s := strconv.FormatInt(prefix, 10)
		full := strings.Repeat(s, repeatCount)
		num, _ := strconv.ParseInt(full, 10, 64)
		ids = append(ids, num)
	}
	return ids
}

func Solve(part2 bool, inputs []string) int64 {
	var invalidIDs []int64
	var wg sync.WaitGroup
	results := make(chan int64, len(inputs))
	for _, input := range inputs {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			continue
		}
		wg.Add(1)
		go func(parts []string) {
			var rangeSum int64 = 0

			defer wg.Done()
			start, _ := strconv.ParseInt(parts[0], 10, 64)
			end, _ := strconv.ParseInt(parts[1], 10, 64)

			startLen := len(parts[0])
			endLen := len(parts[1])
			for length := startLen; length <= endLen; length++ {
				if !part2 {
					if length%2 != 0 {
						continue // only even lengths can be repeated twice
					}
					halfLen := length / 2
					invalidIDs = append(invalidIDs, generateInvalidIDs(length, halfLen)...)
				} else {
					for partLength := 1; partLength <= length/2; partLength++ {
						invalidIDs = append(invalidIDs, generateInvalidIDs(length, partLength)...)
					}
				}

			}
			invalidIDs = removeDuplicate(invalidIDs)
			for _, id := range invalidIDs {
				if id >= start && id <= end {
					// fmt.Println("Invalid ID :", id)
					rangeSum += id
				}
			}
			results <- rangeSum
		}(parts)
	}

	wg.Wait()
	close(results)
	var totalSum int64 = 0
	for sum := range results {
		totalSum += sum
	}
	return totalSum
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
