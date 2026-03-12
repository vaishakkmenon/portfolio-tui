package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	COLS = 28 // Reduced from 40
	ROWS = 14 // Reduced from 20
)

var (
	baseCircuitColor = lipgloss.Color("#2A2D32")
	neonCyan         = lipgloss.Color("#00FFFF")
)

type Point struct{ X, Y int }

type Cell struct {
	pathID int // -1 = empty, 0 = core, 1+ = paths
}

type Grid struct {
	cells [ROWS][COLS]Cell
	Paths [][]Point
}

func NewGrid() *Grid {
	g := &Grid{Paths: [][]Point{}}
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			g.cells[y][x].pathID = -1
		}
	}
	return g
}

func (g *Grid) AddPath(id int, segments []Point) {
	var fullPath []Point
	for i := 0; i < len(segments)-1; i++ {
		curr, next := segments[i], segments[i+1]
		if curr.X == next.X { // Vertical
			step := 1
			if curr.Y > next.Y {
				step = -1
			}
			for y := curr.Y; y != next.Y+step; y += step {
				// Protect the Core (0) from being overwritten
				if g.cells[y][curr.X].pathID != 0 {
					g.cells[y][curr.X].pathID = id
				}
				fullPath = append(fullPath, Point{curr.X, y})
			}
		} else { // Horizontal
			step := 1
			if curr.X > next.X {
				step = -1
			}
			for x := curr.X; x != next.X+step; x += step {
				// Protect the Core (0) from being overwritten
				if g.cells[curr.Y][x].pathID != 0 {
					g.cells[curr.Y][x].pathID = id
				}
				fullPath = append(fullPath, Point{x, curr.Y})
			}
		}
	}
	g.Paths = append(g.Paths, fullPath)
}

// Update the signature to accept the indices instead of a single point
// Render now takes an array tracking the progress of ALL paths
func (g *Grid) Render(pathProgress []int, baseColor lipgloss.Color, fadeColors []lipgloss.Color) string {
	var sb strings.Builder

	trail := make(map[Point]lipgloss.Color)

	// fadeColors := []lipgloss.Color{
	// 	lipgloss.Color("#FF4500"), // Head (Blazing Orange-Red)
	// 	lipgloss.Color("#FF7F50"), // Tail 1 (Coral)
	// 	lipgloss.Color("#CD5C5C"), // Tail 2 (Indian Red)
	// 	lipgloss.Color("#8B0000"), // Tail 3 (Dark Red)
	// 	lipgloss.Color("#3E0000"), // Tail 4 (Deep Rust)
	//}

	// Loop through all paths (skip 0, which is the core)
	for pathIdx := 1; pathIdx < len(g.Paths); pathIdx++ {
		pointIdx := pathProgress[pathIdx]
		if pointIdx < 0 {
			continue // This path is currently idle
		}

		path := g.Paths[pathIdx]
		// Calculate the fading tail for this specific active path
		for i := 0; i < len(fadeColors); i++ {
			pIdx := pointIdx - i
			if pIdx >= 0 && pIdx < len(path) {
				if _, exists := trail[path[pIdx]]; !exists {
					trail[path[pIdx]] = fadeColors[i]
				}
			}
		}
	}

	for y := 0; y < ROWS; y++ {
		var line strings.Builder
		hasContent := false
		for x := 0; x < COLS; x++ {
			myPath := g.cells[y][x].pathID
			if myPath == -1 {
				// Draw the thinner "V" Monogram
				var vChar string

				if (y == 5 && x == 12) || (y == 6 && x == 13) {
					vChar = "\\" // Left arm
				} else if (y == 5 && x == 16) || (y == 6 && x == 15) {
					vChar = "/" // Right arm
				} else if y == 7 && x == 14 {
					vChar = "V" // Bottom point
				}

				if vChar != "" {
					line.WriteString(lipgloss.NewStyle().Foreground(baseColor).Render(vChar))
					continue
				}

				line.WriteRune(' ')
				continue
			}
			hasContent = true

			connects := func(nx, ny int) bool {
				if nx < 0 || nx >= COLS || ny < 0 || ny >= ROWS {
					return false
				}
				neighborPath := g.cells[ny][nx].pathID
				if neighborPath == -1 {
					return false
				}
				return neighborPath == myPath || neighborPath == 0 || myPath == 0
			}

			u, d, l, r := connects(x, y-1), connects(x, y+1), connects(x-1, y), connects(x+1, y)

			var char rune
			switch {
			case u && d && l && r:
				char = '╋'
			case u && d && l:
				char = '┫'
			case u && d && r:
				char = '┣'
			case l && r && u:
				char = '┻'
			case l && r && d:
				char = '┳'
			case u && l:
				char = '┛'
			case u && r:
				char = '┗'
			case d && l:
				char = '┓'
			case d && r:
				char = '┏'
			case u || d:
				char = '┃'
			case l || r:
				char = '━'
			default:
				char = '·'
			}

			color := baseCircuitColor
			if trailColor, ok := trail[Point{x, y}]; ok {
				color = trailColor
			}

			line.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(char)))
		}
		if hasContent {
			sb.WriteString(strings.TrimRight(line.String(), " ") + "\n")
		}
	}
	return sb.String()
}

func BuildChip() *Grid {
	g := NewGrid()

	// Core (CPU Box) - Squished slightly to be 4 units high instead of 6
	g.AddPath(0, []Point{{9, 4}, {19, 4}, {19, 8}, {9, 8}, {9, 4}})

	// Top paths (1-3)
	// Notice they now start at y=1, leaving y=0 empty as a buffer from the shell prompt
	g.AddPath(1, []Point{{10, 1}, {10, 2}, {11, 2}, {11, 4}})
	g.AddPath(2, []Point{{15, 1}, {15, 2}, {14, 2}, {14, 4}})
	g.AddPath(3, []Point{{22, 1}, {18, 1}, {18, 2}, {17, 2}, {17, 4}})

	// Left paths (4-6)
	g.AddPath(4, []Point{{2, 3}, {6, 3}, {6, 5}, {9, 5}})
	g.AddPath(5, []Point{{4, 6}, {9, 6}})
	g.AddPath(6, []Point{{1, 9}, {5, 9}, {5, 7}, {9, 7}})

	// Right paths (7-9)
	g.AddPath(7, []Point{{25, 3}, {22, 3}, {22, 5}, {19, 5}})
	g.AddPath(8, []Point{{26, 6}, {19, 6}})
	g.AddPath(9, []Point{{24, 9}, {21, 9}, {21, 7}, {19, 7}})

	// Bottom paths (10-12)
	// Raised the bottom pins so they don't drag down as far
	g.AddPath(10, []Point{{8, 12}, {8, 10}, {11, 10}, {11, 8}})
	g.AddPath(11, []Point{{13, 12}, {13, 10}, {14, 10}, {14, 8}})
	g.AddPath(12, []Point{{20, 12}, {20, 10}, {17, 10}, {17, 8}})

	return g
}
