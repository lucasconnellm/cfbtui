package teams

import (
	"context"
	"log"
	"net/http"

	"github.com/antihax/optional"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	"github.com/lucasconnellm/cfbtui/pkg/keybindings"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type Props struct {
	Conference optional.String
	SetTeam    func(team gocfbd.Team)

	Client *cfbd.CfbdClient
	Ctx    context.Context
}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Style       lipgloss.Style
	Sort        optional.String
	Teams       []gocfbd.Team
	Table       table.Model
	KeyBindings keybindings.KeyMap
}

func New(ctx context.Context, client *cfbd.CfbdClient) *Component {
	teamTable := table.New(table.WithFocused(true), table.WithWidth(80), table.WithColumns([]table.Column{{Title: "Name", Width: 50}, {Title: "Mascot", Width: 30}}))

	return &Component{
		Style:       lipgloss.Style{},
		Sort:        optional.String{},
		KeyBindings: keybindings.DefaultKeyMap(),
		Table:       teamTable,
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	client := c.Props().Client
	ctx := c.Props().Ctx
	resp, httpResp, err := client.SwaggerClient.TeamsApi.GetTeams(ctx, &gocfbd.TeamsApiGetTeamsOpts{
		Conference: c.Props().Conference,
	})
	if err != nil {
		log.Fatalf("Error retrieving teams: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code %d", httpResp.StatusCode)
	}

	c.Teams = resp

	items := make([]table.Row, 0)
	for _, team := range resp {
		items = append(items, table.Row{team.School, team.Mascot})
	}
	c.Table.SetRows(items)

	c.Table.Focus()
	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.KeyBindings.Quit):
			return tea.Quit
		case key.Matches(msg, c.KeyBindings.Select):
			team := c.Teams[c.Table.Cursor()]
			c.Props().SetTeam(team)
			reactea.SetCurrentRoute("team")
		}
	}
	update, cmd := c.Table.Update(msg)
	c.Table = update
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (c *Component) Render(int, int) string {
	return c.Table.View()
}
