package chat

import (
	"clingy-client/http3"
	"clingy-client/pages/contact"
	"clingy-client/shared"
	"encoding/json"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// TODO: change this to only be called when not already polling
	cmds = append(cmds, WaitForMessage())

	switch msg := msg.(type) {
	case IncomingMessageMsg:
		m.messages = append(m.messages, msg.Message)
		m.setChatContent()
		m.viewport.GotoBottom()

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
				if m.currentChat == nil {
					break
				}
				message := http3.ChatMessage{
					To:      m.currentChat.ID,
					From:    m.config.UniqueID,
					Message: m.chatInput.Value(),
				}
				bytes, _ := json.Marshal(message)
				http3.SendMessage(bytes)
				m.messages = append(m.messages, message)
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
	case contact.SelectChatMsg:
		m.currentChat = &msg.Chat
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
