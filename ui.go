package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"
)

// --- Configuration ---
const (
	lockedHeight  = 18 // Fixed height to prevent vertical jitter
	contentWidth  = 60 // Fixed width for the information pane
	pulseInterval = 30 * time.Millisecond
)

// --- Theme Architecture ---
type Theme struct {
	Brand    lipgloss.Color
	Subtitle lipgloss.Color
	Base     lipgloss.Color
	Fades    []lipgloss.Color
}

var themes = map[string]Theme{
	"FORGE": {
		Brand:    "#FF4500",
		Subtitle: "#CD5C5C",
		Base:     "#1C1C1C",
		Fades:    []lipgloss.Color{"#FF4500", "#FF7F50", "#CD5C5C", "#8B0000", "#3E0000"},
	},
	"NEON": {
		Brand:    "#00FFFF",
		Subtitle: "#00BFFF",
		Base:     "#0A2463",
		Fades:    []lipgloss.Color{"#00FFFF", "#00BFFF", "#1E90FF", "#0000CD", "#000080"},
	},
	"PULSE": {
		Brand:    "#FF00FF",
		Subtitle: "#C71585",
		Base:     "#4A0000",
		Fades:    []lipgloss.Color{"#FF00FF", "#C71585", "#8B008B", "#4B0082", "#2F0000"},
	},
	"TERMINAL": {
		Brand:    "#00FF41",
		Subtitle: "#008F11",
		Base:     "#0D1117",
		Fades:    []lipgloss.Color{"#00FF41", "#008F11", "#003B00", "#002500", "#001000"},
	},
	"DRIFT": {
		Brand:    "#BD93F9",
		Subtitle: "#6272A4",
		Base:     "#282A36",
		Fades:    []lipgloss.Color{"#BD93F9", "#FF79C6", "#8BE9FD", "#50FA7B", "#F1FA8C"},
	},
}

// --- Model & State ---
type model struct {
	choices        []string
	cursor         int
	selected       string
	chipGrid       *Grid
	pathProgress   []int
	moveTicker     int
	currentTheme   Theme
	terminalWidth  int
	terminalHeight int
	sess           ssh.Session
	showHelp       bool
}

type pulseMsg struct{}

func doPulse() tea.Cmd {
	return tea.Tick(pulseInterval, func(t time.Time) tea.Msg {
		return pulseMsg{}
	})
}

func initialModel(s ssh.Session) model {
	grid := BuildChip()
	progress := make([]int, len(grid.Paths))
	for i := range progress {
		progress[i] = -1
	}

	return model{
		sess:         s,
		choices:      []string{"About Me", "Projects", "Certifications", "Contact", "Help"},
		selected:     "About Me",
		chipGrid:     grid,
		pathProgress: progress,
		currentTheme: themes["NEON"],
		showHelp:     true, // Trigger the "First Run" help modal
	}
}

// --- Bubble Tea Lifecycle ---
func (m model) Init() tea.Cmd {
	return doPulse()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth, m.terminalHeight = msg.Width, msg.Height

	case pulseMsg:
		m.moveTicker++
		if m.moveTicker >= 2 {
			m.moveTicker = 0
			for i := 1; i < len(m.pathProgress); i++ {
				if m.pathProgress[i] >= 0 {
					m.pathProgress[i]++
					if m.pathProgress[i] >= len(m.chipGrid.Paths[i])+5 {
						m.pathProgress[i] = -1
					}
				} else if rand.Intn(100) < 2 {
					m.pathProgress[i] = 0
				}
			}
		}
		return m, doPulse()

	case tea.KeyMsg:
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "left", "a", "up", "w":
			m.cursor = (m.cursor - 1 + len(m.choices)) % len(m.choices)
		case "right", "d", "down", "s":
			m.cursor = (m.cursor + 1) % len(m.choices)
		case "h", "?":
			m.showHelp = true
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

// --- View Rendering ---
func (m model) View() string {
	renderer := lipgloss.NewRenderer(m.sess)
	renderer.SetColorProfile(termenv.TrueColor)

	// Header Styles
	brandStyle := renderer.NewStyle().Foreground(m.currentTheme.Brand).Bold(true)
	subStyle := renderer.NewStyle().Foreground(m.currentTheme.Subtitle).Faint(true)
	activeTab := brandStyle.Copy().Underline(true)
	inactiveTab := renderer.NewStyle().Foreground(lipgloss.Color("#666666"))

	// 1. Navigation Row
	var tabs []string
	for i, choice := range m.choices {
		style := inactiveTab
		if i == m.cursor {
			style = activeTab
		}

		// If it's the first tab, remove the left padding to align with your name
		if i == 0 {
			tabs = append(tabs, style.Padding(0, 1, 0, 0).Render(choice))
		} else {
			// Give other tabs some breathing room on both sides
			tabs = append(tabs, style.Padding(0, 1).Render(choice))
		}
	}

	header := renderer.NewStyle().MarginBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))

	// 2. Right Pane (Fixed Content Box)
	contentBody := renderer.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(contentWidth).
		Height(12).
		Render(m.getContent())

	rightPane := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		brandStyle.Render("VAISHAK MENON"),
		lipgloss.JoinHorizontal(lipgloss.Top, subStyle.Italic(true).Render("Software Engineer | Business Analyst")),
		"\n"+brandStyle.Render("── "+m.selected+" ──"),
		"\n"+contentBody,
	)

	// 3. Left Pane (Centered Chip with Divider)
	chipView := m.chipGrid.Render(renderer, m.pathProgress, m.currentTheme.Base, m.currentTheme.Fades)
	topPad := (lockedHeight - lipgloss.Height(chipView)) / 2

	leftPane := renderer.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(topPad+1, 2, 0, 0).
		Height(lockedHeight).
		Render(chipView)

	// 4. Assemble & Help Overlay
	mainUI := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, renderer.NewStyle().PaddingLeft(2).Render(rightPane))

	if m.showHelp {
		return m.renderHelpModal(renderer)
	}

	hPad := max(0, (m.terminalWidth-lipgloss.Width(mainUI))/2)
	vPad := max(0, (m.terminalHeight-lipgloss.Height(mainUI))/2)

	return renderer.NewStyle().Padding(vPad, 0, 0, hPad).Render(mainUI)
}

// --- Sub-View Helpers ---
func (m model) renderHelpModal(r *lipgloss.Renderer) string {
	header := r.NewStyle().Foreground(m.currentTheme.Brand).Bold(true).Render("HOW TO NAVIGATE")
	content := r.NewStyle().Foreground(lipgloss.Color("252")).Render(
		"←/→ or A/D : Switch Tabs\n" +
			"1 - 5       : Change Themes\n" +
			"Q / ESC     : Exit Portfolio\n\n" +
			"Press any key to start...",
	)

	modal := r.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(m.currentTheme.Brand).
		Padding(1, 2).
		Width(40).
		Render(lipgloss.JoinVertical(lipgloss.Center, header, "\n", content))

	return lipgloss.Place(m.terminalWidth, m.terminalHeight, lipgloss.Center, lipgloss.Center, modal)
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
	case "Help":
		return "── NAVIGATION ──\n" +
			"• ← / →  : Switch sections\n" +
			"• A / D  : Switch sections\n\n" +
			"── THEMES ──\n" +
			"• 1 - 5  : Change color palette\n\n" +
			"── SYSTEM ──\n" +
			"• Q / ESC: Exit portfolio\n" +
			"• Ctrl+C : Force quit"
	}
	return ""
}
