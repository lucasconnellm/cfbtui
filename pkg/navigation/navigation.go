package navigation

import (
	"github.com/antihax/optional"
	gocfbd "github.com/lucasconnellm/gocfbd"
)

type NavigateTo interface{}

type NavigateToTeams struct {
	NavigateTo
	Conference optional.String
}

type NavigateToTeam struct {
	NavigateTo
	Team gocfbd.Team
}

type Navigation interface {
	NavigateTo | NavigateToTeam | NavigateToTeams
}
