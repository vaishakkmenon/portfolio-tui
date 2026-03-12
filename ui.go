package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

// --- Model Definition ---

type model struct {
	choices        []string
	cursor         int
	selected       string
	chipGrid       *Grid
	pathProgress   []int // NEW: Tracks the pointIdx for EVERY path independently
	moveTicker     int
	currentTheme   Theme
	terminalWidth  int
	terminalHeight int
	sess           ssh.Session
}

type pulseMsg struct{}

type Theme struct {
	Brand    lipgloss.Color
	Subtitle lipgloss.Color
	Base     lipgloss.Color
	Fades    []lipgloss.Color
}

var themes = map[string]Theme{
	"FORGE": { // Formerly Vantage
		Brand:    lipgloss.Color("#FF4500"),
		Subtitle: lipgloss.Color("#CD5C5C"),
		Base:     lipgloss.Color("#1C1C1C"),
		Fades:    []lipgloss.Color{"#FF4500", "#FF7F50", "#CD5C5C", "#8B0000", "#3E0000"},
	},
	"NEON": { // Formerly Hydrodash
		Brand:    lipgloss.Color("#00FFFF"),
		Subtitle: lipgloss.Color("#00BFFF"),
		Base:     lipgloss.Color("#0A2463"),
		Fades:    []lipgloss.Color{"#00FFFF", "#00BFFF", "#1E90FF", "#0000CD", "#000080"},
	},
	"PULSE": { // Formerly Domain
		Brand:    lipgloss.Color("#FF00FF"),
		Subtitle: lipgloss.Color("#C71585"),
		Base:     lipgloss.Color("#4A0000"),
		Fades:    []lipgloss.Color{"#FF00FF", "#C71585", "#8B008B", "#4B0082", "#2F0000"},
	},
	"TERMINAL": { // Formerly Kube
		Brand:    lipgloss.Color("#00FF41"),
		Subtitle: lipgloss.Color("#008F11"),
		Base:     lipgloss.Color("#0D1117"),
		Fades:    []lipgloss.Color{"#00FF41", "#008F11", "#003B00", "#002500", "#001000"},
	},
	"DRIFT": { // Formerly Crosstrek
		Brand:    lipgloss.Color("#BD93F9"),
		Subtitle: lipgloss.Color("#6272A4"),
		Base:     lipgloss.Color("#282A36"),
		Fades:    []lipgloss.Color{"#BD93F9", "#FF79C6", "#8BE9FD", "#50FA7B", "#F1FA8C"},
	},
}

// --- Helper Functions ---
func doPulse() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return pulseMsg{}
	})
}

// --- Bubble Tea Lifecycle ---
func initialModel(s ssh.Session) model {
	// Initialize all paths as idle (-1)
	grid := BuildChip()
	progress := make([]int, len(grid.Paths))
	for i := range progress {
		progress[i] = -1
	}

	return model{
		sess:         s,
		choices:      []string{"About Me", "Projects", "Certifications", "Contact"},
		selected:     "About Me",
		chipGrid:     grid,
		pathProgress: progress,
		moveTicker:   0,
		currentTheme: themes["NEON"],
	}
}

func (m model) Init() tea.Cmd {
	return doPulse()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		return m, nil
	case pulseMsg:
		// 1. Handle the Cascading Grid Movement
		m.moveTicker++
		if m.moveTicker >= 2 { // Speed of the data streams
			m.moveTicker = 0

			// Loop through every path (skip 0, which is the core)
			for i := 1; i < len(m.pathProgress); i++ {
				if m.pathProgress[i] >= 0 {
					// If the path is active, move the light forward
					m.pathProgress[i]++
					pathLen := len(m.chipGrid.Paths[i])
					trailLen := 5

					// Once the entire tail drains into the core, set to idle
					if m.pathProgress[i] >= pathLen+trailLen {
						m.pathProgress[i] = -1
					}
				} else {
					// If the path is idle, give it a random chance to fire
					// A 2% chance per tick keeps the cascade looking natural and busy
					if rand.Intn(100) < 2 {
						m.pathProgress[i] = 0
					}
				}
			}
		}
		return m, doPulse()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "a", "up", "w":
			m.cursor = (m.cursor - 1 + len(m.choices)) % len(m.choices)
		case "right", "d", "down", "s":
			m.cursor = (m.cursor + 1) % len(m.choices)
		case "1":
			m.currentTheme = themes["FORGE"]
		case "2":
			m.currentTheme = themes["NEON"]
		case "3":
			m.currentTheme = themes["PULSE"]
		case "4":
			m.currentTheme = themes["TERMINAL"]
		case "5":
			m.currentTheme = themes["DRIFT"]
		}
		m.selected = m.choices[m.cursor]
	}
	return m, nil
}

func (m model) View() string {
	renderer := lipgloss.NewRenderer(m.sess)

	// RE-DEFINE TABS USING THE RENDERER
	activeTab := renderer.NewStyle().
		Foreground(m.currentTheme.Brand).
		Underline(true).
		Bold(true)

	inactiveTab := renderer.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	brandColor := m.currentTheme.Brand
	subtitleColor := m.currentTheme.Subtitle

	// 1. Set a STANDARDIZED height for the entire UI
	// 18-20 rows is usually perfect for your current content
	const staticHeight = 16

	// 2. Build the Tabs
	var tabs []string
	for i, choice := range m.choices {
		var tStyle lipgloss.Style
		if i == m.cursor {
			tStyle = activeTab
		} else {
			tStyle = inactiveTab
		}

		if i == 0 {
			tabs = append(tabs, tStyle.Padding(0, 1, 0, 0).Render(choice))
		} else {
			tabs = append(tabs, tStyle.Padding(0, 1).Render(choice))
		}
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	var activeName string
	for name, t := range themes {
		if t.Brand == m.currentTheme.Brand {
			activeName = name
		}
	}

	statusIndicator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Faint(true).
		Render(fmt.Sprintf(" [MODE: %s]", activeName))

	fullHeader := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		headerRow,
		lipgloss.NewStyle().PaddingLeft(10).Render(statusIndicator),
	)

	paddedHeader := lipgloss.NewStyle().MarginBottom(1).Render(fullHeader)

	// 3. Build the Right Pane
	nameStyle := renderer.NewStyle().Foreground(brandColor).Bold(true).SetString("VAISHAK MENON")
	roleStyle := renderer.NewStyle().Foreground(subtitleColor).Faint(true).Italic(true)
	sectionHeaderStyle := renderer.NewStyle().Foreground(brandColor).Bold(true)
	contentBodyStyle := renderer.NewStyle().Foreground(lipgloss.Color("252")).Width(60)
	dividerStyle := renderer.NewStyle().Foreground(subtitleColor).Faint(true)

	// Join them horizontally
	roleLine := lipgloss.JoinHorizontal(lipgloss.Top,
		roleStyle.Render("Software Engineer"),
		dividerStyle.Render(" | "),
		roleStyle.Render("Business Analyst"),
	)

	// Use roleLine in your rightPane stack
	rightPane := lipgloss.JoinVertical(
		lipgloss.Left,
		paddedHeader,
		nameStyle.Render(),
		roleLine, // Use the new horizontally joined line here
		"\n"+sectionHeaderStyle.Render("── "+m.selected+" ──"),
		"\n"+contentBodyStyle.Render(m.getContent()),
	)

	// 4. Render the Chip and calculate the fixed centering
	chipView := m.chipGrid.Render(m.pathProgress, m.currentTheme.Base, m.currentTheme.Fades)
	chipHeight := lipgloss.Height(chipView)

	// Calculate vertical center
	topPad := (staticHeight - chipHeight) / 2

	// NUDGE OFFSET: Change this to +1 to move it up slightly,
	// or set to 0 for perfect mathematical centering
	topPad += 1

	if topPad < 0 {
		topPad = 0
	}

	// 5. Build the Left Pane
	leftPane := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("#333333")).
		PaddingTop(topPad).
		PaddingRight(2).
		Height(staticHeight).
		Render(chipView)

	var mainUI string

	// Check if the screen is too narrow for a side-by-side layout
	if m.terminalWidth < 90 {
		// MOBILE/VERTICAL MODE: Stack them!
		// We remove the border from the leftPane so it doesn't look weird when stacked
		mobileLeftPane := lipgloss.NewStyle().
			PaddingBottom(1).
			Render(m.chipGrid.Render(m.pathProgress, m.currentTheme.Base, m.currentTheme.Fades))

		mainUI = lipgloss.JoinVertical(lipgloss.Left, mobileLeftPane, rightPane)
	} else {
		// DESKTOP MODE: Original side-by-side layout
		mainUI = lipgloss.JoinHorizontal(lipgloss.Top, leftPane, lipgloss.NewStyle().PaddingLeft(2).Render(rightPane))
	}

	uiWidth := lipgloss.Width(mainUI)
	uiHeight := lipgloss.Height(mainUI)

	hPad := (m.terminalWidth - uiWidth) / 2
	vPad := (m.terminalHeight - uiHeight) / 2
	if hPad < 0 {
		hPad = 0
	}
	if vPad < 0 {
		vPad = 0
	}

	return lipgloss.NewStyle().PaddingLeft(hPad).PaddingTop(vPad).Render(mainUI)
}

func (m model) getContent() string {
	switch m.selected {
	case "About Me":
		return "📍 Based in Dillon, MT\n🏢 Business Analyst @ Speridian (New Wave Program)\n☁️  Cloud Native (CKA & AWS Cloud/AI Practitioner)\n\nPassionate about full-stack development, AI/ML, and system architecture. When I'm away from the keyboard, I'm usually analyzing chess positions, researching PC hardware, or meal prepping for the week."
	case "Projects":
		return "• vaishakmenon.com: Personal site.\n• RAG Chatbot: Chatbot about me, my experience, and what I am doing.\n• Vantage: Custom Chess Engine in Rust.\n• Pomodoro Timer: Productivity tool with emphasis with built in sounds/music."
	case "Certifications":
		return "• Certified Kubernetes Administrator (CKA)\n• AWS Certified AI Practitioner\n• AWS Certified Cloud Practitioner"
	case "Contact":
		return "Email: vaishakkmenon25@gmail.com\nGitHub: github.com/vaishakkmenon\nLinkedIn: linkedin.com/in/vaishakkmenon\nCurrently In: Dillon, Montana"
	}
	return ""
}
