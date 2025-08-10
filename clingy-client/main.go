package main

// An example program demonstrating the pager component from the Bubbles
// component library

import (
	"fmt"
	"os"

	"clingy-client/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := app.InitialModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
