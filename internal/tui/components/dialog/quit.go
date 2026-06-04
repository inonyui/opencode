package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/layout"
	"github.com/opencode-ai/opencode/internal/tui/styles"
	"github.com/opencode-ai/opencode/internal/tui/theme"
	"github.com/opencode-ai/opencode/internal/tui/util"
)

const question = "确定要退出吗？"

type CloseQuitMsg struct{}

type QuitDialog interface {
	tea.Model
	layout.Bindings
}

type quitDialogCmp struct {
	selectedNo bool
}

type helpMapping struct {
	LeftRight  key.Binding
	EnterSpace key.Binding
	Yes        key.Binding
	No         key.Binding
	Tab        key.Binding
}

var helpKeys = helpMapping{
	LeftRight: key.NewBinding(
		key.WithKeys("left", "right"),
		key.WithHelp("←/→", "切换选项"),
	),
	EnterSpace: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "确认"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y", "Y"),
		key.WithHelp("y/Y", "是"),
	),
	No: key.NewBinding(
		key.WithKeys("n", "N"),
		key.WithHelp("n/N", "否"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "切换选项"),
	),
}

func (q *quitDialogCmp) Init() tea.Cmd {
	return nil
}

func (q *quitDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpKeys.LeftRight) || key.Matches(msg, helpKeys.Tab):
			q.selectedNo = !q.selectedNo
			return q, nil
		case key.Matches(msg, helpKeys.EnterSpace):
			if !q.selectedNo {
				return q, tea.Quit
			}
			return q, util.CmdHandler(CloseQuitMsg{})
		case key.Matches(msg, helpKeys.Yes):
			return q, tea.Quit
		case key.Matches(msg, helpKeys.No):
			return q, util.CmdHandler(CloseQuitMsg{})
		}
	}
	return q, nil
}

func (q *quitDialogCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()
	
	yesStyle := baseStyle
	noStyle := baseStyle
	spacerStyle := baseStyle.Background(t.Background())

	if q.selectedNo {
		noStyle = noStyle.Background(t.Primary()).Foreground(t.Background())
		yesStyle = yesStyle.Background(t.Background()).Foreground(t.Primary())
	} else {
		yesStyle = yesStyle.Background(t.Primary()).Foreground(t.Background())
		noStyle = noStyle.Background(t.Background()).Foreground(t.Primary())
	}

	yesButton := yesStyle.Padding(0, 1).Render("是")
	noButton := noStyle.Padding(0, 1).Render("否")

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesButton, spacerStyle.Render("  "), noButton)

	width := lipgloss.Width(question)
	remainingWidth := width - lipgloss.Width(buttons)
	if remainingWidth > 0 {
		buttons = spacerStyle.Render(strings.Repeat(" ", remainingWidth)) + buttons
	}

	content := baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			question,
			"",
			buttons,
		),
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (q *quitDialogCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(helpKeys)
}

func NewQuitCmp() QuitDialog {
	return &quitDialogCmp{
		selectedNo: true,
	}
}
