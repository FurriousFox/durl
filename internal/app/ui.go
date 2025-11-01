package ui

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"

	table "argv.nl/durl/internal/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var spinner = []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var spincount = 0

type Model struct {
	Mu     sync.RWMutex
	width  int
	height int
	Focus  int
	Table  table.Model
	IpLen  int
	State  map[string]map[string]any
}

func (m *Model) Init() tea.Cmd {
	// return tick()
	// tea.SetWindowTitle("test")

	maxWidth := len("IP")

	m.Mu.RLock()
	for ip := range m.State {
		if l := len(ip); l > maxWidth {
			maxWidth = l
		}
	}
	m.Mu.RUnlock()

	columns := []table.Column{
		{Title: "IP", Width: maxWidth},
	}

	m.Mu.RLock()
	ips := make([]string, 0, len(m.State))
	for ip := range m.State {
		ips = append(ips, ip)
	}
	sort.Strings(ips)

	rows := []table.Row{}
	for _, ip := range ips {
		rows = append(rows, table.Row{ip})
	}
	m.Mu.RUnlock()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(m.height-3),
		table.WithTruncate(false),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.Mu.Lock()
	m.Focus = 0
	m.Table = t
	m.IpLen = maxWidth
	m.Mu.Unlock()

	// return nil
	return tick()
}

var goodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // green
var badStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))  // red
var spinStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")) // grey

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	maxWidth := len("IP")

	m.Mu.RLock()
	for ip := range m.State {
		if l := len(ip); l > maxWidth {
			maxWidth = l
		}
	}
	m.Mu.RUnlock()

	columns := []table.Column{
		{Title: "IP", Width: maxWidth + 2},
	}

	m.Mu.RLock()
	ips := make([]string, 0, len(m.State))
	for ip := range m.State {
		ips = append(ips, ip)
	}
	sort.Strings(ips)

	rows := []table.Row{}
	for _, ip := range ips {
		value, ok := m.State[ip]["test"].([]any)
		re := regexp.MustCompile(`(\x1b\[[0-9;]*m)+$`)
		if ok && value[0] == true {
			rows = append(rows, table.Row{re.ReplaceAllString(goodStyle.Render(fmt.Sprintf("● %s", ip)), "")})
		} else if ok && value[0] == false {
			rows = append(rows, table.Row{re.ReplaceAllString(badStyle.Render(fmt.Sprintf("● %s", ip)), "")})
		} else {
			rows = append(rows, table.Row{re.ReplaceAllString(spinStyle.Render(fmt.Sprintf("%s %s", spinner[spincount%10], ip)), "")})
		}
		// rows = append(rows, table.Row{ip})
	}
	m.Mu.RUnlock()

	m.Mu.Lock()
	m.IpLen = maxWidth + 2

	m.Table.SetColumns(columns)
	m.Table.SetRows(rows)

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.Table.SetHeight(max(msg.Height-2, 3))

		// case time.Time:
		// 	m--
		// 	if m <= 0 {
		// 		return m, tea.Quit
		// 	}
		// 	return m, tick()
	}
	m.Mu.Unlock()

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
	case time.Time:
		// Just trigger a re-render
		return m, tick()

	}

	m.Mu.Lock()
	m.Table, cmd = m.Table.Update(message)
	m.Mu.Unlock()

	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m *Model) View() string {
	spin := spinner[spincount%10]

	var focusedModelStyle = lipgloss.NewStyle().
		Width(m.width-m.IpLen-6).
		Height(m.height-4).
		Align(lipgloss.Top, lipgloss.Top). /* horizontal, vertical */
		PaddingLeft(1).
		MarginTop(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	if len(m.Table.SelectedRow()) == 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, baseStyle.Render(m.Table.View()), focusedModelStyle.Render(""))
	}

	var stateString string = ""

	m.Mu.RLock()
	ips := make([]string, 0, len(m.State))
	for ip := range m.State {
		ips = append(ips, ip)
	}
	sort.Strings(ips)

	selectedIndex := -1
	selectedRow := m.Table.SelectedRow()
	for i, row := range m.Table.Rows() {
		if row[0] == selectedRow[0] {
			selectedIndex = i
			break
		}
	}
	selectedIP := ips[selectedIndex]

	if m.State[selectedIP]["tcp"] == true {
		stateString += goodStyle.Render("● TCP")
	} else if m.State[selectedIP]["tcp"] != nil {
		stateString += badStyle.Render("● TCP")
	} else {
		stateString += spinStyle.Render(spin + " TCP")
	}

	if m.State[selectedIP]["tls_10"] == true {
		stateString += goodStyle.Render("\n\n● TLS 1.0")
	} else if m.State[selectedIP]["tls_10"] != nil {
		stateString += badStyle.Render("\n\n● TLS 1.0")
	} else {
		stateString += spinStyle.Render("\n\n" + spin + " TLS 1.0")
	}
	if m.State[selectedIP]["tls_11"] == true {
		stateString += goodStyle.Render("\n● TLS 1.1")
	} else if m.State[selectedIP]["tls_11"] != nil {
		stateString += badStyle.Render("\n● TLS 1.1")
	} else {
		stateString += spinStyle.Render("\n" + spin + " TLS 1.1")
	}
	if m.State[selectedIP]["tls_12"] == true {
		stateString += goodStyle.Render("\n● TLS 1.2")
	} else if m.State[selectedIP]["tls_12"] != nil {
		stateString += badStyle.Render("\n● TLS 1.2")
	} else {
		stateString += spinStyle.Render("\n" + spin + " TLS 1.2")
	}
	if m.State[selectedIP]["tls_13"] == true {
		stateString += goodStyle.Render("\n● TLS 1.3")
	} else if m.State[selectedIP]["tls_13"] != nil {
		stateString += badStyle.Render("\n● TLS 1.3")
	} else {
		stateString += spinStyle.Render("\n" + spin + " TLS 1.3")
	}

	value, ok := m.State[selectedIP]["http_11"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n\n● HTTP/1.1 (%s)", value[1]))
	} else if ok {
		stateString += badStyle.Render("\n\n● HTTP/1.1")
	} else {
		stateString += spinStyle.Render("\n\n" + spin + " HTTP/1.1")
	}
	value, ok = m.State[selectedIP]["http_20"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n● HTTP/2   (%s)", value[1]))
	} else if ok {
		stateString += badStyle.Render("\n● HTTP/2")
	} else {
		stateString += spinStyle.Render("\n" + spin + " HTTP/2")
	}
	value, ok = m.State[selectedIP]["http_30"].([]any)
	if ok && value[0] == true {
		stateString += goodStyle.Render(fmt.Sprintf("\n● HTTP/3   (%s)", value[1]))
	} else if ok {
		stateString += badStyle.Render("\n● HTTP/3")
	} else {
		stateString += spinStyle.Render("\n" + spin + " HTTP/3")
	}
	m.Mu.RUnlock()

	return lipgloss.JoinHorizontal(lipgloss.Bottom, baseStyle.Render(m.Table.View()), focusedModelStyle.Render(stateString))
	// return baseStyle.Render(m.table.View())
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		spincount++
		return time.Time(t)
	})
}

func RunUI(model *Model) {
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
