package main

import (
	"clingy-client/app"
	"clingy-client/http3"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	http3.ConnectToServer()
	m := app.InitialModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
