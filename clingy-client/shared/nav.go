package shared

type PageType int

const (
	ChatPage PageType = iota
	ContactPage
	ConnectPage
)

type NavigateMsg struct {
	Page PageType
}
