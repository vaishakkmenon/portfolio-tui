package main

import "github.com/charmbracelet/lipgloss"

var (
	asciiStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Padding(1, 2)

	contentStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)

	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
)
