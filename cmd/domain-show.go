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
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.digTable.SelectedRow()[1]),
			)
		}
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

type Report struct {
	digResults []DNSLookupResult
}

func DomainDetails(host string) {
	// I want to call dig on this and display the information in a table
	// I want to call namecheap and show the same information in a table

	// namecheap.DomainsDNSGetHostsCommandResponse()
	lookupResult := Report{}
	for _, recordType := range []string{"A", "MX", "TXT", "CNAME", "NS"} {
		DNSlookupReports, _ := DigDomain(host, recordType)
		for _, digReport := range *DNSlookupReports {
			lookupResult.digResults = append(lookupResult.digResults, digReport)
		}
	}
	columns := []table.Column{
		{Title: "Type", Width: 8},
		{Title: "Name", Width: 25},
		{Title: "Data", Width: 33},
	}
	rows := []table.Row{}
	for _, result := range lookupResult.digResults {
		fmt.Printf("question: %s\n", result.Question.Type)
		for _, answer := range result.Answer {
			rows = append(rows, table.Row{answer.Type, answer.Name, answer.Data})
		}
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
		table.WithFocused(true),
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
	m := tableModel{t2, t}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
