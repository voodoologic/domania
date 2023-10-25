package dig

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type responseMsg struct {
	quitting bool
}

type responseChannel struct{}

type startModel struct {
	sub       chan struct{}
	responses int
	spinner   spinner.Model
	quitting  bool
}

func callNamecheap() tea.Msg {
	//this is where the activity is called
	time.Sleep(time.Second * 5)
	return responseMsg{
		quitting: true,
	}
}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseChannel(<-sub)
	}
}

func (m startModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		callNamecheap,
		waitForActivity(m.sub),
	)
}

func (m startModel) View() string {
	s := fmt.Sprintf("Fetching your domains from namecheap %s", m.spinner.View())
	if m.quitting {
		s += "\n"
	}
	return s
}

func (m startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case responseMsg:
		response := msg.(responseMsg)
		if response.quitting {
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
	p := tea.NewProgram(startModel{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program", err)
		os.Exit(1)
	}
	items := []list.Item{
		item("dougheadley.com"),
		item("classykathy.com"),
		item("nothingburger.recipies"),
	}
	return items
}
