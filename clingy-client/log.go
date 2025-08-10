package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func log(s string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", s)
}
