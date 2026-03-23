package auth

import "github.com/charmbracelet/lipgloss"

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
