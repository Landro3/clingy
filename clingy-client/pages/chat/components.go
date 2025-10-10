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
)

func (m Model) headerView() string {
	titleText := "Unknown"
	if m.currentChat != nil {
		titleText = m.currentChat.Username
	}
	title := titleStyle.Render(titleText)
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

	chatInputStyle = chatInputStyle.Width(m.viewport.Width - lipgloss.Width(m.contactButtonView()) - lipgloss.Width(m.contactButtonView()) - 12)
	return chatInputStyle.Render(m.chatInput.View())
}

func (m Model) contactButtonView() string {
	if m.focus == Contact {
		return shared.FocusedButtonStyle.Render("Contact")
	}

	return shared.ButtonStyle.Render("Contact")
}

func (m Model) configButtonView() string {
	if m.focus == Config {
		return shared.FocusedButtonStyle.Render("Config")
	}

	return shared.ButtonStyle.Render("Config")
}

func (m Model) formatMessages() string {
	var content strings.Builder
	for _, msg := range m.messages {
		var msgText string
		if msg.From == m.config.UniqueID {
			msgText = fmt.Sprintf("%s [%s]\n%s", msg.From, "TBD:TBD", msg.Message)
		} else {
			msgText = fmt.Sprintf("[%s] %s\n%s", "TBD:TBD", msg.From, msg.Message)
		}

		var line string
		if msg.From == m.config.UniqueID {
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
