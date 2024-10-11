package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	rand "math/rand/v2"
	"os"
	"strings"
	"time"
)

/*
The universe of the Game of Life is an infinite two-dimensional
orthogonal grid of square cells, each of which is in one of two possible states,
live or dead (or populated and unpopulated, respectively).

Every cell interacts with its eight neighbors,
which are the cells that are horizontally, vertically, or diagonally adjacent.

At each step in time, the following transitions occur:

Any live cell with fewer than two live neighbours dies, as if by underpopulation.
Any live cell with two or three live neighbours lives on to the next generation.
Any live cell with more than three live neighbours dies, as if by overpopulation.
Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
The initial pattern constitutes the seed of the system.

The first generation is created by applying the above rules simultaneously to every cell in the seed, live or dead;
births and deaths occur simultaneously, and the discrete moment at which this happens is sometimes called a tick.

Each generation is a pure function of the preceding one. The
rules continue to be applied repeatedly to create further generations.
*/
const tickTime = time.Duration(time.Second / 20)

var gridWidth int
var gridHeight int

type Cell int

const (
	Alive Cell = iota
	Dead
)

func NextGeneration(c Cell, aliveNeighbors int) Cell {
	if c == Alive && (aliveNeighbors < 2 || aliveNeighbors > 3) {
		return Dead

	}
	if c == Dead && aliveNeighbors == 3 {
		return Alive
	}
	return c
}

type Grid struct {
	x     int
	y     int
	Cells [][]Cell
}

func NewGrid(x, y int) *Grid {
	g := Grid{
		x:     x,
		y:     y,
		Cells: make([][]Cell, y),
	}
	for i := range g.Cells {
		g.Cells[i] = make([]Cell, x)
	}
	return &g
}

func (g *Grid) Populate() {
	for y := range g.y {
		line := make([]Cell, g.x)
		if rand.Int64N(3)%3 == 0 {
			for x := range g.x {
				if rand.Int64N(20)%20 == 0 {
					line[x] = Alive
				} else {
					line[x] = Dead
				}
			}
		}

		g.Cells[y] = line
	}
}

func (g *Grid) GetAliveNeighborCells(x, y int) int {
	aliveNeighbors := 0
	n := []int{-1, 0, 1}
	for _, i := range n {
		for _, i2 := range n {
			if i == 0 && i2 == 0 {
				continue
			}
			ny := y + i
			nx := x + i2
			if ny >= 0 && ny < g.y {
				if nx >= 0 && nx < g.x {
					if g.Cells[ny][nx] == Alive {
						aliveNeighbors++
					}
				}
			}
		}
	}
	return aliveNeighbors
}

func (g *Grid) Walk() {
	var buffer bytes.Buffer
	buffer.WriteString(strings.Repeat("\033[1A", g.y+1))
	for y := range g.Cells {
		for x := range g.Cells[y] {
			c := g.Cells[y][x]
			g.Cells[y][x] = NextGeneration(c, g.GetAliveNeighborCells(x, y))
			switch c {
			case Alive:
				buffer.WriteString("█")
			default:
				buffer.WriteString(" ")
			}
		}
	}
	buffer.WriteString("\n")
	fmt.Fprint(os.Stdout, buffer.String())
	buffer.Reset()
}

func main() {
	fd := int(os.Stdout.Fd())
	if !terminal.IsTerminal(fd) {
		fmt.Errorf("not a terminal")
		return
	}
	gridWidth, gridHeight, err := terminal.GetSize(fd)
	gridWidth--
	gridHeight--
	//gridWidth, gridHeight = 300, 60
	if err != nil {
		fmt.Errorf("error getting terminal size %v", err)
		return
	}
	if gridWidth == 0 || gridHeight == 0 {
		fmt.Errorf("could not get terminal size, gridWidth=%d, gridHeight=%d", gridWidth, gridHeight)
		return
	}
	g := NewGrid(gridWidth, gridHeight)
	g.Populate()

	for {
		g.Walk()
		time.Sleep(tickTime)
	}
}
