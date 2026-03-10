package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices    []string
	cursor     int
	selected   string
	blink      bool
	blinkID    int
	logoColors []string
	colorIndex int
}

type blinkMsg struct {
	id int
}

type colorMsg struct{}

func doBlink(id int) tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return blinkMsg{id: id}
	})
}

func doColorTick() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return colorMsg{}
	})
}

func initialModel() model {
	return model{
		choices:    []string{"About Me", "Projects", "Certifications", "Contact"},
		selected:   "About Me",
		blink:      true,
		blinkID:    0,
		logoColors: []string{"#2A0080", "#4B0082", "#6A0DAD", "#8A2BE2", "#9400D3", "#8A2BE2", "#6A0DAD", "#4B0082"},
		colorIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(doBlink(m.blinkID), doColorTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case blinkMsg:
		if msg.id == m.blinkID {
			m.blink = !m.blink
			return m, doBlink(m.blinkID)
		}
		return m, nil

	case colorMsg:
		if m.colorIndex < len(m.logoColors)-1 {
			m.colorIndex++
		} else {
			m.colorIndex = 0
		}
		return m, doColorTick()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "w":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}

			m.blink = true
			m.blinkID++
			return m, doBlink(m.blinkID)

		case "down", "s":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

			m.blink = true
			m.blinkID++
			return m, doBlink(m.blinkID)

		case "enter", " ":
			m.selected = m.choices[m.cursor]
		}
	}
	return m, nil
}

func (m model) View() string {

	currentColor := m.logoColors[m.colorIndex]
	dynamicLogoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(currentColor))

	var menuBuilder strings.Builder
	menuBuilder.WriteString(titleStyle.Render("VAISHAK MENON") + "\n")
	menuBuilder.WriteString("Software Engineer & Business Analyst\n\n")

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			if m.blink {
				cursor = cursorStyle.Render("> ")
			} else {
				cursor = cursorStyle.Render("  ")
			}
			choice = selectedStyle.Render(choice)
		}
		fmt.Fprintf(&menuBuilder, "%s%s\n", cursor, choice)
	}
	menuBuilder.WriteString("\n(Press 'q' to quit)\n")

	var content string
	switch m.selected {
	case "About Me":
		content = "Currently a participant in the New Wave program at\nSperidian. Passionate about AI, Rust, and systems architecture."
	case "Projects":
		content = "• Vantage: A custom Chess Engine built in Rust.\n• vaishakmenon.com: Personal site w/ RAG Chatbot.\n• Pomodoro Timer: Clean, minimal Pomodoro timer with music and ambient sounds."
	case "Certifications":
		content = "• Certified Kubernetes Administrator (CKA)\n• AWS Certified AI Practitioner\n• AWS Certified Cloud Practitioner"
	case "Contact":
		content = "GitHub:   github.com/vaishakkmenon\nLinkedIn: linkedin.com/in/vaishakkmenon\nCurrent Location: Dillon, Montana"
	}

	rightPaneContent := lipgloss.JoinVertical(lipgloss.Left, menuBuilder.String(), "\n"+titleStyle.Render("--- "+m.selected+" ---")+"\n\n"+content)
	leftPane := dynamicLogoStyle.Render(initialANSIShadowASCII)
	rightPane := contentStyle.Render(rightPaneContent)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}
