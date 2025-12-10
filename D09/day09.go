package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct{ x, y int }
type Line struct {
	p1, p2 Point
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func areaWith(rt, other *Point) int {
	return (abs(rt.x-other.x) + 1) * (abs(rt.y-other.y) + 1)
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

func solvePart1(puzzle []string) (result int) {

	points := getRedTiles(puzzle)

	for i, p1 := range points {
		for j, p2 := range points {
			if i != j {
				area := int(math.Abs(float64(p2.x-p1.x+1)) * math.Abs(float64(p2.y-p1.y+1)))
				if area > result {
					result = area
				}
			}
		}
	}

	return
}

func getRedTiles(puzzle []string) []Point {
	var points []Point
	for _, row := range puzzle {
		values := strings.Split(row, ",")
		x, _ := strconv.Atoi(values[0])
		y, _ := strconv.Atoi(values[1])
		points = append(points, Point{x, y})

	}
	return points
}

func getGreenSegments(polygon []Point) []Line {
	segments := make([]Line, 0, len(polygon)+1)
	for i := 0; i < len(polygon)-1; i++ {
		segment := Line{p1: polygon[i], p2: polygon[i+1]}
		segments = append(segments, segment)
	}

	// Complete the polygon
	segments = append(segments, Line{p1: polygon[len(polygon)-1], p2: polygon[0]})
	return segments
}

func (s *Line) intersects(rectA Point, rectB Point) bool {
	recMinX := min(rectA.x, rectB.x) + 1
	recMaxX := max(rectA.x, rectB.x) - 1
	recMinY := min(rectA.y, rectB.y) + 1
	recMaxY := max(rectA.y, rectB.y) - 1

	segMinX := min(s.p1.x, s.p2.x)
	segMaxX := max(s.p1.x, s.p2.x)
	segMinY := min(s.p1.y, s.p2.y)
	segMaxY := max(s.p1.y, s.p2.y)

	if segMaxX < recMinX || segMinX > recMaxX {
		return false
	}
	if segMaxY < recMinY || segMinY > recMaxY {
		return false
	}
	return true
}

func solvePart2(puzzle []string) int {
	points := getRedTiles(puzzle)
	greenSegments := getGreenSegments(points)
	maxRect := 0
	for i := 0; i < len(points)-1; i++ {
	main:
		for j := i + 1; j < len(points); j++ {
			area := areaWith(&points[i], &points[j])
			if area < maxRect {
				continue
			}
			for _, greenSegment := range greenSegments {
				if greenSegment.intersects(points[i], points[j]) {
					continue main
				}
			}

			maxRect = area
		}
	}
	return maxRect
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
	sumPart2 := solvePart2(fileContent)
	fmt.Println("Part 1:", sumPart1)
	fmt.Println("Part 2:", sumPart2)

	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
