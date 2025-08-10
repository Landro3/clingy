package chat

import (
	"fmt"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), m.chatInputView(), m.footerView())
}
