package app

import (
	"context"

	"github.com/antihax/optional"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	"github.com/lucasconnellm/cfbtui/views/team"
	"github.com/lucasconnellm/cfbtui/views/teams"
	gocfbd "github.com/lucasconnellm/gocfbd"
	"github.com/spf13/viper"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	team gocfbd.Team
	game gocfbd.Game
}

func New() *Component {
	return &Component{
		mainRouter: router.New(),
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	backgroundContext := context.Background()
	ctx := context.WithValue(backgroundContext, gocfbd.ContextAccessToken, viper.GetString("cfbd_key"))

	client := cfbd.GetClient()

	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {

			component := teams.New(ctx, client)
			return component, component.Init(teams.Props{
				Conference: optional.NewString("SEC"),
				SetTeam: func(team gocfbd.Team) {
					c.team = team
				},
				Client: client,
				Ctx:    ctx,
			})
		},
		"team": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := team.New(ctx, client)
			return component, component.Init(team.Props{
				Team: c.team,
				SetGame: func(game gocfbd.Game) {
					c.game = game
				},
				Client: client,
				Ctx:    ctx,
			})
		},
	})
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return reactea.Destroy
		}
	}
	return c.mainRouter.Update(msg)
}

func (c *Component) Render(width int, height int) string {
	return c.mainRouter.Render(width, height)
}
