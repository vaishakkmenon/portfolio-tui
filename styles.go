package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8A2BE2"))
	contentStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)

	tabRowStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("240"))

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4500")). // Blazing Orange
			Underline(true).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666")) // Dark, muted grey

	rightPaneStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(60)

	dividerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#333333")).
			PaddingRight(2)
)
