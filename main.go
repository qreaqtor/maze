package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"iter"
	"os"
	"strconv"
	"text/template"
)

const (
	numberRangeErrText = "expected number should be in the range from 0 to 9, but was {{.}}\n"
	numberTypeErrText  = "expected value to be int, but was {{.}}\n"
)

var (
	numberRangeErrTemp, numberTypeErrTemp *template.Template

	inputReader               *bufio.Reader
	outputWriter, errorWriter *bufio.Writer
)

const defaultErrExitCode = 1

type point struct {
	X, Y int
}

func init() {
	numberRangeErrTemp = template.Must(template.New("numberRangeErr").Parse(numberRangeErrText))
	numberTypeErrTemp = template.Must(template.New("numberTypeErr").Parse(numberTypeErrText))

	inputReader = bufio.NewReader(os.Stdin)
	outputWriter = bufio.NewWriter(os.Stdout)
	errorWriter = bufio.NewWriter(os.Stderr)
}

func main() {
	var n, m int
	fmt.Fscan(inputReader, &n, &m)

	maze := readMaze(n, m)

	var start, end point
	fmt.Fscan(inputReader, &start.X, &start.Y, &end.X, &end.Y)

	for current := range findPath(maze, start, end) {
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
				numberTypeErrTemp.Execute(errorWriter, rawNumber)
				os.Exit(defaultErrExitCode)
			}
			if number < 0 || number > 9 {
				numberRangeErrTemp.Execute(errorWriter, number)
				os.Exit(defaultErrExitCode)
			}

			maze[i][j] = int(number)
		}
	}

	return maze
}

func findPath(maze [][]int, start, end point) iter.Seq[point] {
	pq := new(PriorityQueue[point])
	heap.Push(pq, NewItem(start, 0))

	cameFrom := make(map[point]point)
	cameFrom[start] = point{-1, -1}

	costs := make(map[point]int)
	costs[start] = 0

	n, m := len(maze), len(maze[0])

	heap.Init(pq)
	for pq.Len() != 0 {
		item, ok := pq.Pop().(*Item[point])
		if !ok {
			os.Exit(defaultErrExitCode)
		}
		current := item.value

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

			if v, ok := costs[next]; !ok || cost < v {
				costs[next] = cost
				priority := cost + calculateCost(next, end)
				heap.Push(pq, NewItem(next, cost+priority))
				cameFrom[next] = current
			}
		}
	}

	return func(yield func(point) bool) {
		from := end
		for {
			next, ok := cameFrom[from]
			if !(ok && yield(from)) {
				return
			}
			from = next
		}
	}
}

var moves = []point{
	point{-1, 0},
	point{0, -1},
	point{1, 0},
	point{0, 1},
}

func calculateCost(from, to point) int {
	return abs(from.X-to.X) + abs(from.Y-to.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Item[T any] struct {
	value    T   // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
}

func NewItem[T any](value T, priority int) *Item[T] {
	return &Item[T]{
		value:    value,
		priority: priority,
	}
}

type PriorityQueue[T any] []*Item[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}
