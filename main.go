package main

import (
	"bufio"
	"container/heap"
	"fmt"
	priorityqueue "maze/priorityQueue"
	"os"
	"strconv"
	"text/template"
)

const (
	numberRangeErrText = "expected number should be in the range from 0 to 9, but was {{.}}\n"
	numberTypeErrText  = "expected value to be int, but was {{.}}\n"

	noPathText = "no paths found"

	badStartInput = "the start coordinate should not be a wall (0 value in maze)"
	badEndInput   = "the end coordinate should not be a wall (0 value in maze)"
)

var (
	numberRangeErrTemp, numberTypeErrTemp *template.Template

	inputReader  *bufio.Reader
	outputWriter *bufio.Writer
)

const (
	badInputExitCode = 1
	badTypeExitCode  = 2

	noPathExitCode = 3
)

type point struct {
	X, Y int
}

func init() {
	numberRangeErrTemp = template.Must(template.New("numberRangeErr").Parse(numberRangeErrText))
	numberTypeErrTemp = template.Must(template.New("numberTypeErr").Parse(numberTypeErrText))

	inputReader = bufio.NewReader(os.Stdin)
	outputWriter = bufio.NewWriter(os.Stdout)
}

func main() {
	var n, m int
	fmt.Fscan(inputReader, &n, &m)

	maze := readMaze(n, m)

	var start, end point
	fmt.Fscan(inputReader, &start.X, &start.Y, &end.X, &end.Y)

	if maze[start.X][start.Y] == 0 {
		fmt.Fprintln(os.Stderr, badStartInput)
		os.Exit(badInputExitCode)
	} else if maze[end.X][end.Y] == 0 {
		fmt.Fprintln(os.Stderr, badEndInput)
		os.Exit(badInputExitCode)
	}

	for _, current := range findShortestPath(maze, start, end) {
		fmt.Fprintln(outputWriter, current.X, current.Y)
	}
	fmt.Fprintln(outputWriter, ".")

	outputWriter.Flush()
}

func readMaze(n, m int) [][]int {
	maze := make([][]int, n)

	var rawNumber string
	for i := range n {
		maze[i] = make([]int, m)

		for j := range m {
			fmt.Fscan(inputReader, &rawNumber)

			number, err := strconv.Atoi(rawNumber)
			if err != nil {
				numberTypeErrTemp.Execute(os.Stderr, rawNumber)
				os.Exit(badTypeExitCode)
			}
			if number < 0 || number > 9 {
				numberRangeErrTemp.Execute(os.Stderr, number)
				os.Exit(badInputExitCode)
			}

			maze[i][j] = int(number)
		}
	}

	return maze
}

var moves = []point{
	{-1, 0},
	{0, -1},
	{1, 0},
	{0, 1},
}

func findShortestPath(maze [][]int, start, end point) []point {
	minPathLength := calculateDistance(start, end)
	n, m := len(maze), len(maze[0])

	cameFrom := make(map[point]point, minPathLength)

	costs := make(map[point]int, minPathLength)
	costs[start] = 0

	pointsPQ := priorityqueue.NewPriorityQueue[point](minPathLength)
	heap.Push(pointsPQ, priorityqueue.NewItem(start, 0))

	for pointsPQ.Len() != 0 {
		current := heap.Pop(pointsPQ).(*priorityqueue.Item[point]).Value

		if current.X == end.X && current.Y == end.Y {
			break
		}

		for _, next := range moves {
			next.X += current.X
			next.Y += current.Y

			if next.X < 0 || next.X == n || next.Y < 0 || next.Y == m || maze[next.X][next.Y] == 0 {
				continue
			}

			cost := costs[current] + maze[next.X][next.Y]

			if prevCost, ok := costs[next]; !ok || cost < prevCost {
				costs[next] = cost
				priority := cost + calculateDistance(next, end)
				heap.Push(pointsPQ, priorityqueue.NewItem(next, priority))
				cameFrom[next] = current
			}
		}
	}

	if _, ok := cameFrom[end]; !ok {
		fmt.Fprintln(os.Stderr, noPathText)
		os.Exit(noPathExitCode)
	}

	return recoverPath(cameFrom, start, end, 0)
}

func recoverPath(cameFrom map[point]point, start, end point, depth int) []point {
	if start == end {
		path := make([]point, 0, depth+1)
		return append(path, end)
	}

	path := recoverPath(cameFrom, start, cameFrom[end], depth+1)
	return append(path, end)
}

func calculateDistance(from, to point) int {
	return abs(from.X-to.X) + abs(from.Y-to.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
