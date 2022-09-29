package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/antihax/optional"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	"github.com/lucasconnellm/cfbtui/pkg/keybindings"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type Team struct {
	swaggerTeam gocfbd.Team
}

type TeamDelegate struct{}

type Model struct {
	Style      lipgloss.Style
	Sort       optional.String
	Conference optional.String
	KeyMap     keybindings.KeyMap
	Teams      []Team
	List       list.Model
}

func (team Team) FilterValue() string {
	return team.swaggerTeam.School
}

func (del TeamDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	team := item.(Team)
	buf := bytes.NewBufferString(team.swaggerTeam.School)
	w.Write(buf.Bytes())
}

func (del TeamDelegate) Height() int {
	return 1
}

func (del TeamDelegate) Spacing() int {
	return 1
}

func (del TeamDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.List.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	update, cmd := m.List.Update(msg)
	m.List = update
	return m, cmd
}

func (m Model) Init() tea.Cmd {
	return nil
}

type NewOpts struct {
	conference optional.String
}

func New(ctx context.Context, client *cfbd.CfbdClient, opts *NewOpts) Model {
	resp, httpResp, err := client.SwaggerClient.TeamsApi.GetTeams(ctx, &gocfbd.TeamsApiGetTeamsOpts{
		Conference: opts.conference,
	})
	if err != nil {
		log.Fatalf("Error retrieving teams: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code %d", httpResp.StatusCode)
	}
	teams := make([]Team, 0)
	for _, team := range resp {
		teams = append(teams, Team{swaggerTeam: team})
	}
	items := make([]list.Item, 0)
	for _, team := range teams {
		items = append(items, team)
	}
	delegate := TeamDelegate{}
	teamList := list.New(items, delegate, 80, 30)

	return Model{
		Style:      lipgloss.Style{},
		Sort:       optional.String{},
		Conference: optional.String{},
		KeyMap:     keybindings.KeyMap{},
		Teams:      teams,
		List:       teamList,
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
