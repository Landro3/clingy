package chat

import (
	"clingy-client/shared"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	leftChatStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(shared.Purple)).
			Padding(0, 1).
			Margin(0, 0, 1)

	rightChatStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(shared.Cyan)).
			Align(lipgloss.Right).
			Padding(0, 1).
			Margin(0, 0, 1)

	chatInputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(shared.Fg))

	focusedButtonStyle = buttonStyle.BorderForeground(lipgloss.Color(shared.Purple))
)

func (m Model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) chatInputView() string {
	if m.focus != ChatInput {
		chatInputStyle = chatInputStyle.BorderForeground(lipgloss.Color(shared.Fg))
	} else {
		chatInputStyle = chatInputStyle.BorderForeground(lipgloss.Color(shared.Purple))
	}

	return chatInputStyle.Render(m.chatInput.View())
}

func (m Model) contactButtonView() string {
	if m.focus != Contact {
		return buttonStyle.Render("Contact")
	}

	return focusedButtonStyle.Render("Contact")
}

func (m Model) configButtonView() string {
	if m.focus != Config {
		return buttonStyle.Render("Config")
	}

	return focusedButtonStyle.Render("Config")
}

func (m Model) formatMessages() string {
	var content strings.Builder
	for _, msg := range m.messages {
		var msgText string
		if msg.IsOwn {
			msgText = fmt.Sprintf("%s [%s]\n%s", msg.Author, msg.Time, msg.Content)
		} else {
			msgText = fmt.Sprintf("[%s] %s\n%s", msg.Time, msg.Author, msg.Content)
		}

		var line string
		if msg.IsOwn {
			boxedMsg := rightChatStyle.Render(msgText)
			line = lipgloss.PlaceHorizontal(m.viewport.Width-2, lipgloss.Right, boxedMsg)
		} else {
			line = leftChatStyle.Render(msgText)
		}

		content.WriteString(line + "\n")
	}
	return content.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
