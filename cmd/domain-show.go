package dig

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/voodoologic/domania/cmd/namecheapHandler"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	digTable  table.Model
	nameCheap table.Model
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.digTable.Focused() {
				m.digTable.Blur()
			} else {
				m.digTable.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, nil
		}
	}
	if len(m.digTable.Rows()) != len(m.nameCheap.Rows()) {
		rows, cheap := consolidate(m)
		rows, cheap = lookupOutliers(rows, cheap)
		m.digTable.SetRows(rows)
		m.digTable.SetHeight(len(rows))
		m.nameCheap.SetRows(cheap)
		m.nameCheap.SetHeight(len(cheap))
	}
	m.digTable, cmd = m.digTable.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	body := strings.Builder{}
	body.WriteString("wut up?\n")
	pad := lipgloss.NewStyle().Padding(1)
	tables := lipgloss.JoinHorizontal(
		lipgloss.Top,
		pad.Render(m.digTable.View()),
		pad.Render(m.nameCheap.View()),
	)
	return baseStyle.Render(tables) + "\n"
}

func consolidate(myTableModel tableModel) ([]table.Row, []table.Row) {
	digTableModel := myTableModel.digTable.Rows()
	cheapsTableModel := myTableModel.nameCheap.Rows()
	for _, row := range digTableModel {
		cheap := returnMatchForRow(row, cheapsTableModel)
		if cheap != nil {
			row[0] = "Y"
			cheap[0] = "Y"
		} else {
			row[0] = "N"
		}
	}
	return digTableModel, cheapsTableModel
}

func returnMatchForRow(row table.Row, cheaps []table.Row) table.Row {
	for _, cheap := range cheaps {
		if match(row, cheap) {
			return cheap
		}
	}

	return nil
}

type Report struct {
	digResults []DNSLookupResult
}

func match(row, cheap table.Row) bool {
	for i := 1; len(row) > i; i++ {
		if row[i] == cheap[i] {
			continue
		} else {
			return false
		}
	}
	return true
}

func lookupOutliers(rows, cheaps []table.Row) ([]table.Row, []table.Row) {

	for _, cheap := range cheaps {
		if cheap[0] == "?" {
			//extract domain
			newRows, err := DigDomain(cheap[2], cheap[1])
			if err != nil {
			}
			rows = append(rows, newRows...)
		}
	}
	return rows, cheaps
}

func DomainDetails(host string) {
	// I want to call dig on this and display the information in a table
	// I want to call namecheap and show the same information in a table
	// namecheap.DomainsDNSGetHostsCommandResponse()
	rows, err := InitDomain(host)
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "M", Width: 1},
		{Title: "Type", Width: 5},
		{Title: "Name", Width: 25},
		{Title: "Data", Width: 33},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
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
	//
	//
	rows2, _ := namecheapHandler.GetDomainDetails(host)
	t2 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows2),
		table.WithFocused(false),
		table.WithHeight(len(rows2)),
	)
	s2 := table.DefaultStyles()
	s2.Header = s2.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s2.Selected = s2.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t2.SetStyles(s2)
	m := tableModel{t, t2}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
