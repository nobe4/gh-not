package colors

import (
	"github.com/fatih/color"
)

func Colorize(c color.Attribute, s string) string {
	return color.New(c).SprintFunc()(s)
}

func Red(s string) string     { return Colorize(color.FgRed, s) }
func Magenta(s string) string { return Colorize(color.FgMagenta, s) }
func Blue(s string) string    { return Colorize(color.FgBlue, s) }
func Cyan(s string) string    { return Colorize(color.FgCyan, s) }
func Green(s string) string   { return Colorize(color.FgGreen, s) }
func Yellow(s string) string  { return Colorize(color.FgYellow, s) }
