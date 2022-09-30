package navigation

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type NavigationStack struct {
	Stack []interface{}
	Index int
}

func PushToStack(nav Navigation) tea.Cmd {
	switch nav := nav.(type) {
	case NavigateToTeam:
		log.Printf("navigate to team: %s", nav.Team.School)
		return nil
	case NavigateToTeams:
		log.Printf("navigate to teams: %s", nav.Conference.Value())
		return nil
	default:
		return nil
	}
}
