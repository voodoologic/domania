package dig

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type responseMsg struct {
	domains  []list.Item
	quitting bool
}

type responseChannel struct {
	domains []list.Item
}

type startModel struct {
	sub       chan responseChannel
	responses int
	spinner   spinner.Model
	quitting  bool
	domains   []list.Item
}

func callNamecheap(sub chan responseChannel) tea.Cmd {
	//this is where the activity is called
	return func() tea.Msg {
		for {
			domains := GetDomains()
			sub <- responseChannel{domains: domains}
		}
	}
}

func waitForActivity(sub chan responseChannel) tea.Cmd {
	return func() tea.Msg {
		return responseChannel(<-sub)
	}
}

func (m startModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		callNamecheap(m.sub),
		waitForActivity(m.sub),
	)
}

func (m startModel) View() string {
	s := fmt.Sprintf("Fetching your domains from namecheap %s", m.spinner.View())
	if m.quitting {
		// s += "\n"
	}
	return s
}

func (m startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case responseChannel:
		response := msg.(responseChannel)
		if response.domains != nil {
			m.domains = response.domains
			return m, tea.Quit
		} else {
			return m, waitForActivity(m.sub)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func StartProgram() []list.Item {
	// start spinner and get domains
	model := startModel{
		sub:     make(chan responseChannel),
		spinner: spinner.New(),
	}
	p := tea.NewProgram(model)
	returnModel, err := p.Run()
	if err != nil {
		fmt.Println("could not start program", err)
		os.Exit(1)
	}
	return returnModel.(startModel).domains
}
