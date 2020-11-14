package authentication

import (
	"context"
	"fmt"
	"log"
	"os"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

//DDApiKey name of dd api env var
const DDApiKey string = "DD_CLIENT_API_KEY"

//DDAppKey name of dd app env var
const DDAppKey string = "DD_CLIENT_APP_KEY"

/*
Authenticate - authenticate against DD API service
*/
func Authenticate() (context.Context, *datadog.APIClient) {

	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: os.Getenv(DDApiKey),
			},
			"appKeyAuth": {
				Key: os.Getenv(DDAppKey),
			},
		},
	)

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	response, httpResponse, err := apiClient.AuthenticationApi.Validate(ctx).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when validating the DD auth keys: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", httpResponse)
		log.Fatal("ERROR") //clean this up
	}

	fmt.Fprintf(os.Stdout, "Validated (response): %v\n", response)

	return ctx, apiClient

}
