package teams

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/antihax/optional"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	"github.com/lucasconnellm/cfbtui/pkg/keybindings"
	"github.com/lucasconnellm/cfbtui/pkg/navigation"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type Model struct {
	Style       lipgloss.Style
	Sort        optional.String
	Conference  optional.String
	Teams       []gocfbd.Team
	Table       table.Model
	KeyBindings keybindings.KeyMap
}

func (m Model) View() string {
	return m.Table.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyBindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyBindings.Select):
			_, cmd := m.ViewTeam()
			cmds = append(cmds, cmd)
		}
	}
	update, cmd := m.Table.Update(msg)
	m.Table = update
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) ViewTeam() (tea.Model, tea.Cmd) {
	team := m.Teams[m.Table.Cursor()]
	return m, func() tea.Msg {
		return navigation.NavigateToTeam{
			Team: team,
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type NewOpts struct {
	conference optional.String
}

func New(ctx context.Context, client *cfbd.CfbdClient, opts *NewOpts) Model {
	resp, httpResp, err := client.SwaggerClient.TeamsApi.GetTeams(ctx, &gocfbd.TeamsApiGetTeamsOpts{
		// Conference: opts.conference,
		Conference: optional.NewString("SEC"),
	})
	if err != nil {
		log.Fatalf("Error retrieving teams: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code %d", httpResp.StatusCode)
	}
	items := make([]table.Row, 0)
	for _, team := range resp {
		items = append(items, table.Row{team.School, team.Mascot})
	}
	teamTable := table.New(table.WithFocused(true), table.WithWidth(80), table.WithColumns([]table.Column{{Title: "Name", Width: 50}, {Title: "Mascot", Width: 30}}), table.WithRows(items))

	return Model{
		Style:       lipgloss.Style{},
		Sort:        optional.String{},
		Conference:  optional.String{},
		Teams:       resp,
		KeyBindings: keybindings.DefaultKeyMap(),
		Table:       teamTable,
	}
}

func main() {
	godotenv.Load(".env")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	backgroundContext := context.Background()
	ctx := context.WithValue(backgroundContext, gocfbd.ContextAccessToken, cfbd.GetKey())
	client := cfbd.GetClient()

	m := New(ctx, client, &NewOpts{})
	log.Println(m.Conference)
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
