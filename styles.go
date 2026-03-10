package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8A2BE2"))
	contentStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)

	tabRowStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("240"))

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A2BE2")).
			Bold(true).
			Underline(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))

	rightPaneStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(60)
)
