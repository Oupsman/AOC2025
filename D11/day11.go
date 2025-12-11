package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type PathState struct {
	node    string
	visited map[string]bool
}

type Graph map[string][]string

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

func parseInput(puzzle []string) Graph {
	graph := make(Graph)
	for _, line := range puzzle {
		nodes := strings.Split(line, ":")
		childs := strings.Fields(nodes[1])
		graph[nodes[0]] = childs
	}
	return graph
}

func countPaths(
	graph map[string][]string,
	start, end string,
	memo map[string]int,
) int {
	if val, ok := memo[start]; ok {
		return val
	}

	if start == end {
		return 1
	}

	count := 0
	for _, neighbor := range graph[start] {
		count += countPaths(graph, neighbor, end, memo)
	}

	memo[start] = count
	return count
}

func countPathsWithConstraints(
	graph map[string][]string,
	start, end string,
	requiredNodes map[string]int,
	visitedMask int,
	memo map[string]map[int]int,
) int {
	// Initialise la map de mémoïsation pour ce nœud si nécessaire
	if memo[start] == nil {
		memo[start] = make(map[int]int)
	}
	if val, ok := memo[start][visitedMask]; ok {
		return val
	}

	// Met à jour le masque pour le nœud actuel (si c'est un nœud obligatoire)
	newMask := visitedMask
	if idx, ok := requiredNodes[start]; ok {
		newMask |= (1 << idx)
	}

	// Si on est à la fin, vérifie si tous les nœuds requis sont visités
	if start == end {
		if newMask == (1<<len(requiredNodes))-1 {
			return 1
		} else {
			return 0
		}
	}

	count := 0
	for _, neighbor := range graph[start] {
		count += countPathsWithConstraints(graph, neighbor, end, requiredNodes, newMask, memo)
	}

	if memo[start] == nil {
		memo[start] = make(map[int]int)
	}
	memo[start][newMask] = count
	return count
}
func solvePart1(puzzle []string) int {
	graph := parseInput(puzzle)

	return countPaths(graph, "you", "out", make(map[string]int))
}

func solvePart2(puzzle []string) int {
	graph := parseInput(puzzle)
	requiredNodes := map[string]int{
		"dac": 0,
		"fft": 1,
	}
	memo := make(map[string]map[int]int)

	nbPaths := countPathsWithConstraints(
		graph, "svr", "out", requiredNodes, 0, memo,
	)

	return nbPaths
}

func main() {
	timeStart := time.Now()
	// INPUT := "sample2.txt"
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
