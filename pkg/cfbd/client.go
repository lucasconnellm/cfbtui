package cfbd

import (
	"net/http"

	swagger "github.com/lucasconnellm/gocfbd"
)

type CfbdClient struct {
	HttpClient    *http.Client
	SwaggerClient *swagger.APIClient
}

func GetClient() *CfbdClient {
	httpClient := &http.Client{}
	conf := swagger.NewConfiguration()
	conf.BasePath = "https://api.collegefootballdata.com"
	conf.HTTPClient = httpClient
	return &CfbdClient{
		HttpClient:    httpClient,
		SwaggerClient: swagger.NewAPIClient(conf),
	}
}
