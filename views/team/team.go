package team

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	Team    gocfbd.Team
	SetGame func(game gocfbd.Game)

	Ctx    context.Context
	Client *cfbd.CfbdClient
}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Style       lipgloss.Style
	Game        []gocfbd.Game
	Table       table.Model
	KeyBindings keybindings.KeyMap
}

func New(ctx context.Context, client *cfbd.CfbdClient) *Component {
	teamTable := table.New(table.WithFocused(true), table.WithWidth(70), table.WithColumns([]table.Column{
		{Title: "Week", Width: 5},
		{Title: "Time", Width: 20},
		{Title: "Opponent", Width: 20},
		{Title: "R", Width: 5},
		{Title: "Score", Width: 20},
	}))

	return &Component{
		Style:       lipgloss.Style{},
		KeyBindings: keybindings.DefaultKeyMap(),
		Table:       teamTable,
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	client := c.Props().Client
	ctx := c.Props().Ctx
	resp, httpResp, err := client.SwaggerClient.GamesApi.GetGames(ctx, 2022, &gocfbd.GamesApiGetGamesOpts{
		Team: optional.NewString(c.Props().Team.School),
	})
	if err != nil {
		log.Fatalf("Error retrieving games: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code %d", httpResp.StatusCode)
	}
	items := make([]table.Row, 0)
	for _, game := range resp {
		isAway := game.AwayTeam == c.Props().Team.School
		var intro string
		if isAway {
			intro = "@"
		} else {
			intro = "vs"
		}

		var opponent string
		if isAway {
			opponent = game.HomeTeam
		} else {
			opponent = game.AwayTeam
		}

		future := false
		var humanStartDate string
		if game.StartTimeTbd {
			future = true
		} else {
			startTime, err := time.Parse(time.RFC3339, game.StartDate)
			humanStartDate = startTime.Format("Jan 02 T 03:04 pm")
			if err != nil {
				log.Fatalf("error parsing timestamp: %v", err)
			}
			if time.Now().Before(startTime) {
				future = true
			}
		}

		var teamPoints, opponentPoints int
		var dubCol string
		var pointsStr string
		if !future {
			if isAway {
				teamPoints = int(game.AwayPoints)
				opponentPoints = int(game.HomePoints)
			} else {
				teamPoints = int(game.HomePoints)
				opponentPoints = int(game.AwayPoints)
			}

			if teamPoints > opponentPoints {
				dubCol = "W"
			} else {
				dubCol = "L"
			}
			pointsStr = fmt.Sprintf("%d-%d", teamPoints, opponentPoints)
		}

		items = append(items, table.Row{
			fmt.Sprintf("%d", game.Week),
			humanStartDate,
			fmt.Sprintf("%s %s", intro, opponent),
			dubCol,
			pointsStr,
		})
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
			log.Printf("no path to view game yet")
			return nil
		case key.Matches(msg, c.KeyBindings.Back):
			reactea.SetCurrentRoute("teams")
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
