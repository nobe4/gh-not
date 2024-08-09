package colors

import (
	"github.com/charmbracelet/lipgloss"
)

func Colorize(c string, s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(s)
}

func Red(s string) string     { return Colorize("1", s) }
func Magenta(s string) string { return Colorize("5", s) }
func Blue(s string) string    { return Colorize("4", s) }
func Cyan(s string) string    { return Colorize("6", s) }
func Green(s string) string   { return Colorize("2", s) }
func Yellow(s string) string  { return Colorize("3", s) }
