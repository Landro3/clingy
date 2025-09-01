package chat

import (
	"clingy-client/http3"
	"clingy-client/shared"
	"clingy-client/util"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case IncomingMessageMsg:
		chatMsg := ChatMessage{
			Author:  msg.Message.Author,
			Content: msg.Message.Content,
			Time:    msg.Message.Time,
			IsOwn:   false,
		}
		m.messages = append(m.messages, chatMsg)
		m.setChatContent()
		m.viewport.GotoBottom()

		cmds = append(cmds, WaitForMessage())

	case tea.KeyMsg:
		if msg.String() == "tab" {
			m.focus += 1

			if m.focus == Max {
				m.focus = 0
			}

			if m.focus == ChatInput {
				m.chatInput.Focus()
			} else {
				m.chatInput.Blur()
			}
		}

		if msg.String() == "enter" {
			switch m.focus {
			case ChatInput:
				util.Log("...sending message")
				util.Log(fmt.Sprintf("messages length: %d", len(m.messages)))
				http3.SendMessage(m.chatInput.Value())
				m.messages = append(m.messages, ChatMessage{Author: "Andrew", Content: m.chatInput.Value(), Time: "now", IsOwn: true})
				m.setChatContent()
				m.chatInput.SetValue("")
				m.viewport.GotoBottom()
			case Contact:
				cmds = append(cmds, shared.NavigateCmd(shared.ContactPage))
			case Config:
				cmds = append(cmds, shared.NavigateCmd(shared.ConfigPage))
			}
		}

		// TODO: reconcile this with viewport update below
		m.chatInput, cmd = m.chatInput.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		chatInputHeight := lipgloss.Height(m.chatInputView())

		m.resizeChatInput(msg.Width)
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-chatInputHeight-verticalMarginHeight)
			m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
			m.viewport.YPosition = headerHeight

			m.ready = true
		} else {
			// Viewport sizing
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight - chatInputHeight
		}

		m.setChatContent()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) resizeChatInput(msgWidth int) {
	if msgWidth < 100 {
		m.chatInput.Width = msgWidth - 12
	} else {
		m.chatInput.Width = msgWidth / 2
	}
}

func (m *Model) setChatContent() {
	m.viewport.SetContent(m.formatMessages())
}
