package main

import (
	"bufio"
	"os/exec"

	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Lamps []int

func applySwitchLamps(state Lamps, switchIndex int, switches [][]int) Lamps {
	newState := make(Lamps, len(state))
	copy(newState, state)
	for _, lamp := range switches[switchIndex] {
		if lamp < len(newState) {
			newState[lamp] = 1 - newState[lamp]
		}
	}
	return newState
}

func lampsToKey(state Lamps) string {
	return fmt.Sprintf("%v", state)
}

func solveLampPuzzle(numLamps int, switches [][]int, target Lamps) ([]int, int) {
	initial := make(Lamps, numLamps)
	visited := make(map[string]bool)
	queue := []struct {
		state Lamps
		path  []int
	}{{initial, []int{}}}
	visited[lampsToKey(initial)] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if lampsToKey(current.state) == lampsToKey(target) {
			return current.path, len(current.path)
		}

		for i := 0; i < len(switches); i++ {
			newState := applySwitchLamps(current.state, i, switches)
			key := lampsToKey(newState)
			if !visited[key] {
				visited[key] = true
				newPath := append(current.path, i)
				queue = append(queue, struct {
					state Lamps
					path  []int
				}{newState, newPath})
			}
		}
	}

	return nil, math.MaxInt
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

func solvePart1(puzzle []string) int {
	sum := 0
	for _, line := range puzzle {
		var interruptors [][]int
		fields := strings.Fields(line)
		targetLampsT := strings.Trim(fields[0], "[]")
		interruptorsValues := fields[1 : len(fields)-1]
		// joltage := fields[len(fields)-1]

		for _, interruptorsT := range interruptorsValues {
			interruptorsT = strings.Trim(interruptorsT, "()")
			interruptorsV := strings.Split(interruptorsT, ",")
			inter := make([]int, 0, len(interruptorsV))
			for _, v := range interruptorsV {
				v = strings.TrimSpace(v) // Nettoie les espaces
				i, err := strconv.Atoi(v)
				if err != nil {
					fmt.Printf("Erreur de conversion pour '%s'\n", v)
					continue
				}
				inter = append(inter, i)
			}
			interruptors = append(interruptors, inter)
		}
		targetLamps := make(Lamps, len(targetLampsT))
		for i, lamp := range targetLampsT {
			if lamp == '#' {
				targetLamps[i] = 1
			}
		}

		_, res := solveLampPuzzle(len(targetLamps), interruptors, targetLamps)
		fmt.Println(targetLamps, ":", res, "(", interruptors, ")")
		sum += res

	}

	return sum
}

type Puzzle struct {
	Interruptors [][]int
	Joltages     []int
}

func parsePuzzleLine(line string) (*Puzzle, error) {
	fields := strings.Fields(line)
	interruptorsValues := fields[1 : len(fields)-1]
	joltagesT := strings.Trim(fields[len(fields)-1], "{}")

	var interruptors [][]int
	for _, interruptorsT := range interruptorsValues {
		interruptorsT = strings.Trim(interruptorsT, "()")
		interruptorsV := strings.Split(interruptorsT, ",")
		var inter []int
		for _, v := range interruptorsV {
			v = strings.TrimSpace(v)
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("erreur de conversion pour '%s': %v", v, err)
			}
			inter = append(inter, i)
		}
		interruptors = append(interruptors, inter)
	}

	joltagesV := strings.Split(joltagesT, ",")
	var joltages []int
	for _, joltT := range joltagesV {
		joltT = strings.TrimSpace(joltT)
		v, err := strconv.Atoi(joltT)
		if err != nil {
			return nil, fmt.Errorf("erreur de conversion pour '%s': %v", joltT, err)
		}
		joltages = append(joltages, v)
	}

	return &Puzzle{Interruptors: interruptors, Joltages: joltages}, nil
}

func generateLPFile(puzzle *Puzzle, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("impossible de créer %s: %v", filename, err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Minimize\n")
	fmt.Fprintf(file, "obj: ")
	for i := range puzzle.Interruptors {
		fmt.Fprintf(file, " + x%d", i+1)
	}
	fmt.Fprintf(file, "\n\nSubject To\n")

	for j := 0; j < len(puzzle.Joltages); j++ {
		fmt.Fprintf(file, "c%d: ", j+1)
		first := true
		for i, inter := range puzzle.Interruptors {
			if contains(inter, j) {
				if !first {
					fmt.Fprintf(file, " + ")
				}
				fmt.Fprintf(file, "x%d", i+1)
				first = false
			}
		}
		fmt.Fprintf(file, " = %d\n", puzzle.Joltages[j])
	}

	fmt.Fprintf(file, "\nBounds\n")
	for i := range puzzle.Interruptors {
		fmt.Fprintf(file, "x%d >= 0\n", i+1)
	}

	fmt.Fprintf(file, "\nGenerals\n")
	for i := range puzzle.Interruptors {
		fmt.Fprintf(file, "x%d\n", i+1)
	}

	fmt.Fprintf(file, "End\n")
	return nil
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func solveWithGLPK(lpFile string) (int, error) {
	cmd := exec.Command("glpsol", "--lp", lpFile, "--output", "solution.txt")
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("GLPK failed: %v", err)
	}

	return parseGLPKObjective("solution.txt")
}

func parseGLPKObjective(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("unable to read %s: %v", filename, err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Objective:") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				var objective int
				_, err := fmt.Sscanf(parts[3], "%d", &objective)
				if err != nil {
					return 0, fmt.Errorf("Unable to parse objective: %v", err)
				}
				return objective, nil
			}
		}
	}
	return 0, fmt.Errorf("Objective not found in %s", filename)
}

func solvePuzzle(puzzle *Puzzle) (int, error) {

	err := generateLPFile(puzzle, "joltages.lp")
	if err != nil {
		return 0, fmt.Errorf("Unable to create LP file: %v", err)
	}

	return solveWithGLPK("joltages.lp")
}

func solvePart2(puzzle []string) int {
	var sum int

	for _, line := range puzzle {
		puzzle, err := parsePuzzleLine(line)
		if err != nil {
			continue
		}

		solution, err := solvePuzzle(puzzle)
		if err == nil {
			sum += solution
		}
	}

	return sum
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
	fmt.Println("================================")
	sumPart2 := solvePart2(fileContent)
	fmt.Println("Part 1:", sumPart1)
	fmt.Println("Part 2:", sumPart2)

	fmt.Printf("Time: %.2fms\n", float64(time.Since(timeStart).Microseconds())/1000)
}
