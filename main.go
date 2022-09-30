package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/antihax/optional"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	"github.com/lucasconnellm/cfbtui/pkg/navigation"
	"github.com/lucasconnellm/cfbtui/views/teams"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type model struct {
	Plays    []gocfbd.Play
	Index    int32
	viewport viewport.Model
	Teams    teams.Model
	NavStack navigation.NavigationStack
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

const width = 78

func GetPlays() []gocfbd.Play {
	backgroundContext := context.Background()
	ctx := context.WithValue(backgroundContext, gocfbd.ContextAccessToken, cfbd.GetKey())

	client := cfbd.GetClient()

	resp, httpResp, _ := client.SwaggerClient.PlaysApi.GetPlays(ctx, 2022, 1, &gocfbd.PlaysApiGetPlaysOpts{Team: optional.NewString("Georgia")})
	log.Println(resp)
	log.Println(httpResp)

	return resp
}

func initialModel() model {
	return model{
		Plays: make([]gocfbd.Play, 0),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case navigation.NavigateToTeam, navigation.NavigateToTeams:
		navigation.PushToStack(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.Plays = GetPlays()

		case "left":
			if m.Index == 0 {
				break
			}
			m.Index--

		case "right":
			m.Index++
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		log.Println(msg)

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "r":
			m.Plays = GetPlays()

		case "left":
			if m.Index == 0 {
				break
			}
			m.Index--

		case "right":
			m.Index++
		}

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func BallOn(ydLine int, viewSize int, pad int) string {
	// ydSize = (viewSize - pad) / 100
	fieldString := ""
	yards := 100
	for i := 0; i < yards; i++ {
		if ydLine == i {
			fieldString += "x"
		} else if i%5 == 0 {
			fieldString += "|"
		} else {
			fieldString += " "
		}
	}
	return fieldString
}

func (m model) View() string {
	content := "No plays. Press \"r\" to refresh"
	if len(m.Plays) != 0 {
		content, _ = getVPContent(fmt.Sprintf("%s\n%s", m.Plays[m.Index].PlayText, BallOn(int(m.Plays[m.Index].YardLine), 200, 10)))
	}
	m.viewport.SetContent(content)
	return m.viewport.View() + m.helpView()
	if len(m.Plays) == 0 {
		return "No plays. Press \"r\" to refresh"
	}
	return fmt.Sprintf("%s\n%s", m.Plays[m.Index].PlayText, BallOn(int(m.Plays[m.Index].YardLine), 200, 10))
}

func (m model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}

const content = `
# Today’s Menu
## Appetizers
| Name        | Price | Notes                           |
| ---         | ---   | ---                             |
| Tsukemono   | $2    | Just an appetizer               |
| Tomato Soup | $4    | Made with San Marzano tomatoes  |
| Okonomiyaki | $4    | Takes a few minutes to make     |
| Curry       | $3    | We can add squash if you’d like |
## Seasonal Dishes
| Name                 | Price | Notes              |
| ---                  | ---   | ---                |
| Steamed bitter melon | $2    | Not so bitter      |
| Takoyaki             | $3    | Fun to eat         |
| Winter squash        | $3    | Today it's pumpkin |
## Desserts
| Name         | Price | Notes                 |
| ---          | ---   | ---                   |
| Dorayaki     | $4    | Looks good on rabbits |
| Banana Split | $5    | A classic             |
| Cream Puff   | $3    | Pretty creamy!        |
All our dishes are made in-house by Karen, our chef. Most of our ingredients
are from our garden or the fish market down the street.
Some famous people that have eaten here lately:
* [x] René Redzepi
* [x] David Chang
* [ ] Jiro Ono (maybe some day)
Bon appétit!
`

func getVPContent(content string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return "", err
	}
	return str, nil

}

func example() (*model, error) {
	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	str, err := getVPContent(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &model{
		Plays:    []gocfbd.Play{},
		Index:    0,
		viewport: vp,
	}, nil

}

func main() {
	godotenv.Load(".env")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m, err := example()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
