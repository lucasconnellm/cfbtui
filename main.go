package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/antihax/optional"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/lucasconnellm/cfbtui/pkg/cfbd"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type model struct {
	Plays []gocfbd.Play
	Index int32
}

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
			m.Index--

		case "right":
			m.Index++
		}

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	if len(m.Plays) == 0 {
		return "No plays. Press \"r\" to refresh"
	}
	return m.Plays[m.Index].PlayText
}

func main() {
	godotenv.Load(".env")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
