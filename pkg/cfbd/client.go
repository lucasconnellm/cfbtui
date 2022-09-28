package cfbd

import (
	"net/http"
	"os"

	swagger "github.com/lucasconnellm/gocfbd"
)

type CfbdClient struct {
	HttpClient    *http.Client
	apikey        string
	SwaggerClient *swagger.APIClient
}

func GetKey() string {
	return os.Getenv("CFBD_KEY")
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
