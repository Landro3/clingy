package util

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func Log(s string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", s)
}
