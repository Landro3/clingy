package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	// Background Colors
	BgDark      = "#1a1b26"
	Bg          = "#24283b"
	BgHighlight = "#414868"
	Comment     = "#565f89"

	// Foreground Colors
	Fg       = "#c0caf5"
	FgDark   = "#a9b1d6"
	FgGutter = "#9aa5ce"

	// Accent Colors
	Blue   = "#7aa2f7"
	Purple = "#bb9af7"
	Cyan   = "#7dcfff"
	Teal   = "#73daca"
	Green  = "#9ece6a"
	Yellow = "#e0af68"
	Red    = "#f7768e"
	Orange = "#ff9e64"
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
			BorderForeground(lipgloss.Color(Purple)).
			Padding(0, 1)

	rightChatStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(Cyan)).
			Align(lipgloss.Right).
			Padding(0, 1)

	chatInputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	contactButtonStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder())
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
		chatInputStyle = chatInputStyle.BorderForeground(lipgloss.Color(Fg))
	} else {
		chatInputStyle = chatInputStyle.BorderForeground(lipgloss.Color(Purple))
	}

	if m.focus != Contact {
		contactButtonStyle = contactButtonStyle.BorderForeground(lipgloss.Color(Fg))
	} else {
		contactButtonStyle = contactButtonStyle.BorderForeground(lipgloss.Color(Purple))
	}

	input := chatInputStyle.Render(m.chatInput.View())

	return lipgloss.PlaceHorizontal(
		m.viewport.Width,
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, input, contactButtonStyle.Render("Contact")),
	)
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
