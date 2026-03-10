package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices      []string
	cursor       int
	selected     string
	blink        bool
	blinkID      int
	fullGradient []string
	colorIndex   int
	direction    int
}

type blinkMsg struct {
	id int
}

type pulseMsg struct{}

func parseHex(s string) ([3]uint8, error) {
	var res [3]uint8
	v, _ := strconv.ParseInt(s[1:], 16, 32)
	res[0] = uint8(v >> 16)
	res[1] = uint8(v >> 8 & 0xFF)
	res[2] = uint8(v & 0xFF)
	return res, nil
}

func generateSmoothGradient(keyframes []string, stepsPerSegment int) []string {
	var fullGradient []string

	for i := 0; i < len(keyframes)-1; i++ {
		start, _ := parseHex(keyframes[i])
		end, _ := parseHex(keyframes[i+1])

		for j := 0; j < stepsPerSegment; j++ {
			ratio := float64(j) / float64(stepsPerSegment)
			r := uint8(float64(start[0]) + ratio*float64(int(end[0])-int(start[0])))
			g := uint8(float64(start[1]) + ratio*float64(int(end[1])-int(start[1])))
			b := uint8(float64(start[2]) + ratio*float64(int(end[2])-int(start[2])))
			fullGradient = append(fullGradient, fmt.Sprintf("#%02x%02x%02x", r, g, b))
		}
	}

	return fullGradient
}

func doBlink(id int) tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return blinkMsg{id: id}
	})
}

func doPulse() tea.Cmd {
	return tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
		return pulseMsg{}
	})
}

func initialModel() model {
	keyframes := []string{
		"#2A0080", "#4B0082", "#6A0DAD", "#8A2BE2",
		"#9400D3", "#8A2BE2", "#6A0DAD", "#4B0082",
	}

	smoothGradient := generateSmoothGradient(keyframes, 20)

	return model{
		choices:      []string{"About Me", "Projects", "Certifications", "Contact"},
		selected:     "About Me",
		blink:        true,
		blinkID:      0,
		fullGradient: smoothGradient,
		colorIndex:   0,
		direction:    1,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(doBlink(m.blinkID), doPulse())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case blinkMsg:
		if msg.id == m.blinkID {
			m.blink = !m.blink
			return m, doBlink(m.blinkID)
		}
		return m, nil

	case pulseMsg:
		m.colorIndex += m.direction
		if m.colorIndex >= len(m.fullGradient)-1 {
			m.colorIndex = len(m.fullGradient) - 1
			m.direction = -1
		} else if m.colorIndex <= 0 {
			m.colorIndex = 0
			m.direction = 1
		}
		return m, doPulse()

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

	currentColor := m.fullGradient[m.colorIndex]
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
		content = "Currently a participant in the New Wave program at\nSperidian. Passionate about Software Engineering, AI, Chess, and Dr.Pepper."
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
