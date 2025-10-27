package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	Focus  int
	Table  table.Model
	IpLen  int
	State  map[string]map[string]any
}

func (m Model) Init() tea.Cmd {
	// return tick()
	// tea.SetWindowTitle("test")

	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.Table.SetHeight(msg.Height - 2)

		// case time.Time:
		// 	m--
		// 	if m <= 0 {
		// 		return m, tea.Quit
		// 	}
		// 	return m, tick()
	}

	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "esc":
		// 	if m.table.Focused() {
		// 		m.table.Blur()
		// 	} else {
		// 		m.table.Focus()
		// 	}
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
			// case "enter":
			// 	return m, tea.Batch(
			// 		tea.Printf("Let's go to %s!", m.table.SelectedRow()[0]),
			// 	)
		}
	}
	m.Table, cmd = m.Table.Update(message)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() string {
	var focusedModelStyle = lipgloss.NewStyle().
		Width(m.width-m.IpLen-6).
		Height(m.height-4).
		Align(lipgloss.Top, lipgloss.Top). /* horizontal, vertical */
		PaddingLeft(1).
		MarginTop(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	// currently selected
	goodStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // green
	badStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))  // red

	var stateString string = ""
	if m.State[m.Table.SelectedRow()[0]]["tcp"] == true {
		stateString += goodStyle.Render("● TCP")
	} else {
		stateString += badStyle.Render("● TCP")
	}

	if m.State[m.Table.SelectedRow()[0]]["tls_10"] == true {
		stateString += goodStyle.Render("\n\n● TLS 1.0")
	} else {
		stateString += badStyle.Render("\n\n● TLS 1.0")
	}
	if m.State[m.Table.SelectedRow()[0]]["tls_11"] == true {
		stateString += goodStyle.Render("\n● TLS 1.1")
	} else {
		stateString += badStyle.Render("\n● TLS 1.1")
	}
	if m.State[m.Table.SelectedRow()[0]]["tls_12"] == true {
		stateString += goodStyle.Render("\n● TLS 1.2")
	} else {
		stateString += badStyle.Render("\n● TLS 1.2")
	}
	if m.State[m.Table.SelectedRow()[0]]["tls_13"] == true {
		stateString += goodStyle.Render("\n● TLS 1.3")
	} else {
		stateString += badStyle.Render("\n● TLS 1.3")
	}

	value, ok := m.State[m.Table.SelectedRow()[0]]["http_11"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n\n● HTTP/1.1 (%s)", value[1]))
	} else {
		stateString += badStyle.Render("\n\n● HTTP/1.1")
	}
	value, ok = m.State[m.Table.SelectedRow()[0]]["http_20"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n● HTTP/2   (%s)", value[1]))
	} else {
		stateString += badStyle.Render("\n● HTTP/2")
	}
	value, ok = m.State[m.Table.SelectedRow()[0]]["http_30"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n● HTTP/3   (%s)", value[1]))
	} else {
		stateString += badStyle.Render("\n● HTTP/3")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, baseStyle.Render(m.Table.View()), focusedModelStyle.Render(stateString))
	// return baseStyle.Render(m.table.View())
}

// func tick() tea.Cmd {
// 	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
// 		return time.Time(t)
// 	})
// }
