package cfbd

import (
	"context"
	"log"
	"net/http"

	"github.com/antihax/optional"
	swagger "github.com/lucasconnellm/gocfbd"
)

func (client *CfbdClient) GetScoreboard(ctx context.Context, conference string) []swagger.ScoreboardGame {
	cfbdResp, httpResp, err := client.SwaggerClient.GamesApi.GetScoreboard(ctx, &swagger.GamesApiGetScoreboardOpts{
		Conference: optional.NewString(conference),
	})
	if err != nil {
		log.Fatalln(err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Non 200 code: %d", httpResp.StatusCode)
	}
	return cfbdResp
}
