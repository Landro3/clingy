package chat

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
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
			BorderForeground(lipgloss.Color("#7dcfff")).
			Padding(0, 1)

	rightChatStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#bb9af7")).
			Align(lipgloss.Right).
			Padding(0, 1)
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
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	return lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, style.Render(m.chatInput.View()))
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
