package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette for the chip diagram
var (
	// Box style
	boxColor = lipgloss.Color("#FFFFFF")

	// Path colors - each path gets a unique color
	pathColors = []lipgloss.Color{
		"#FF6B6B", // Path 1  - red
		"#4ECDC4", // Path 2  - teal
		"#45B7D1", // Path 3  - blue
		"#96CEB4", // Path 4  - sage
		"#FFEAA7", // Path 5  - yellow
		"#DDA0DD", // Path 6  - plum
		"#98D8C8", // Path 7  - mint
		"#F7DC6F", // Path 8  - gold
		"#BB8FCE", // Path 9  - purple
		"#F0B27A", // Path 10 - orange
		"#82E0AA", // Path 11 - green
		"#F1948A", // Path 12 - salmon
	}

	// Endpoint circle color
	endpointColor = lipgloss.Color("#FF4757")

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			MarginBottom(1)

	// Container style
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	// Legend style
	legendStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1)
)

const (
	COLS = 48
	ROWS = 28
)

type Cell struct {
	char   rune
	pathID int // -1 = none, 0 = box, 1-12 = path number
}

type Grid struct {
	cells       [ROWS][COLS]Cell
	connections [ROWS][COLS]map[string]bool
	circles     map[[2]int]bool
}

func NewGrid() *Grid {
	g := &Grid{
		circles: make(map[[2]int]bool),
	}
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			g.cells[y][x] = Cell{char: ' ', pathID: -1}
			g.connections[y][x] = make(map[string]bool)
		}
	}
	return g
}

func (g *Grid) connect(gx, gy int, direction string, pathID int) {
	if gx >= 0 && gx < COLS && gy >= 0 && gy < ROWS {
		g.connections[gy][gx][direction] = true
		current := g.cells[gy][gx].pathID
		// Paths (1-12) always override box (0) or unset (-1)
		if current == -1 || (pathID > 0 && current == 0) {
			g.cells[gy][gx].pathID = pathID
		}
	}
}

func (g *Grid) drawHLine(x1, x2, y, pathID int) {
	minX, maxX := x1, x2
	if x1 > x2 {
		minX, maxX = x2, x1
	}
	for gx := minX; gx <= maxX; gx++ {
		if gx > minX {
			g.connect(gx, y, "left", pathID)
		}
		if gx < maxX {
			g.connect(gx, y, "right", pathID)
		}
	}
}

func (g *Grid) drawVLine(x, y1, y2, pathID int) {
	minY, maxY := y1, y2
	if y1 > y2 {
		minY, maxY = y2, y1
	}
	for gy := minY; gy <= maxY; gy++ {
		if gy > minY {
			g.connect(x, gy, "up", pathID)
		}
		if gy < maxY {
			g.connect(x, gy, "down", pathID)
		}
	}
}

func (g *Grid) drawPath(pathID int, segments [][2]int) {
	for i := 0; i < len(segments)-1; i++ {
		x1, y1 := segments[i][0], segments[i][1]
		x2, y2 := segments[i+1][0], segments[i+1][1]
		if x1 == x2 {
			g.drawVLine(x1, y1, y2, pathID)
		} else if y1 == y2 {
			g.drawHLine(x1, x2, y1, pathID)
		}
	}
}

func (g *Grid) drawRect(x, y, w, h int) {
	g.drawHLine(x, x+w, y, 0)
	g.drawHLine(x, x+w, y+h, 0)
	g.drawVLine(x, y, y+h, 0)
	g.drawVLine(x+w, y, y+h, 0)
}

func (g *Grid) markCircle(x, y int) {
	g.circles[[2]int{x, y}] = true
}

func (g *Grid) render() string {
	// Resolve characters
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			if g.circles[[2]int{x, y}] {
				g.cells[y][x].char = 'o'
				continue
			}
			c := g.connections[y][x]
			if len(c) == 0 {
				continue
			}
			hasH := c["left"] || c["right"]
			hasV := c["up"] || c["down"]
			if hasH && hasV {
				g.cells[y][x].char = '+'
			} else if hasH {
				g.cells[y][x].char = '-'
			} else if hasV {
				g.cells[y][x].char = '|'
			}
		}
	}

	// Build styled output
	var sb strings.Builder
	for y := 0; y < ROWS; y++ {
		hasContent := false
		for x := 0; x < COLS; x++ {
			if g.cells[y][x].char != ' ' {
				hasContent = true
				break
			}
		}
		if !hasContent && y < 4 {
			continue
		}

		for x := 0; x < COLS; x++ {
			cell := g.cells[y][x]
			if cell.char == ' ' {
				sb.WriteRune(' ')
				continue
			}

			var style lipgloss.Style
			switch {
			case cell.char == 'o':
				style = lipgloss.NewStyle().
					Foreground(endpointColor).
					Bold(true)
			case cell.pathID == 0:
				style = lipgloss.NewStyle().
					Foreground(boxColor).
					Bold(true)
			case cell.pathID >= 1 && cell.pathID <= 12:
				style = lipgloss.NewStyle().
					Foreground(pathColors[cell.pathID-1])
			default:
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#666666"))
			}
			sb.WriteString(style.Render(string(cell.char)))
		}

		// Trim trailing spaces
		line := strings.TrimRight(sb.String(), " ")
		sb.Reset()
		sb.WriteString(line)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func buildChip() *Grid {
	g := NewGrid()

	// Box
	BL, BR, BT, BB := 18, 30, 10, 16
	g.drawRect(BL, BT, BR-BL, BB-BT)

	// Top paths
	g.drawPath(1, [][2]int{{21, BT}, {21, 7}, {19, 7}, {19, 4}})
	g.markCircle(19, 4)

	g.drawPath(2, [][2]int{{24, BT}, {24, 7}, {26, 7}, {26, 5}})
	g.markCircle(26, 5)

	g.drawPath(3, [][2]int{{27, BT}, {27, 8}, {30, 8}, {30, 6}, {35, 6}})
	g.markCircle(35, 6)

	// Left paths
	// Path 4: left, up 2, left
	g.drawPath(4, [][2]int{{BL, 11}, {14, 11}, {14, 9}, {10, 9}})
	g.markCircle(10, 9)

	// Path 5: short straight left
	g.drawPath(5, [][2]int{{BL, 13}, {14, 13}})
	g.markCircle(14, 13)

	// Path 6: zigzag - left past P5, up 2, left to endpoint
	g.drawPath(6, [][2]int{{BL, 15}, {12, 15}, {12, 13}, {8, 13}})
	g.markCircle(8, 13)

	// Right paths
	g.drawPath(7, [][2]int{{BR, 11}, {35, 11}, {35, 9}, {39, 9}})
	g.markCircle(39, 9)

	g.drawPath(8, [][2]int{{BR, 13}, {34, 13}, {34, 15}, {39, 15}})
	g.markCircle(39, 15)

	g.drawPath(9, [][2]int{{BR, 15}, {32, 15}, {32, 17}, {37, 17}})
	g.markCircle(37, 17)

	// Bottom paths
	g.drawPath(10, [][2]int{{21, BB}, {21, 18}, {19, 18}, {19, 20}, {16, 20}, {16, 22}})
	g.markCircle(16, 22)

	g.drawPath(11, [][2]int{{24, BB}, {24, 19}, {22, 19}, {22, 21}})
	g.markCircle(22, 21)

	g.drawPath(12, [][2]int{{27, BB}, {27, 18}, {29, 18}, {29, 20}, {32, 20}, {32, 22}})
	g.markCircle(32, 22)

	return g
}

func renderLegend() string {
	var parts []string
	labels := []string{
		"P1", "P2", "P3", "P4", "P5", "P6",
		"P7", "P8", "P9", "P10", "P11", "P12",
	}
	for i, label := range labels {
		style := lipgloss.NewStyle().
			Foreground(pathColors[i]).
			Bold(true)
		parts = append(parts, style.Render(label))
	}
	return legendStyle.Render("Paths: " + strings.Join(parts, " "))
}

func main() {
	grid := buildChip()
	chip := grid.render()

	title := titleStyle.Render("⚡ Chip Diagram")
	legend := renderLegend()

	content := chip + legend
	output := containerStyle.Render(content)

	fmt.Println()
	fmt.Println(title)
	fmt.Println(output)
}
