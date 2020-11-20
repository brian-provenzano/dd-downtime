package authentication

import (
	"context"
	"fmt"
	"log"
	"os"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

//DDApiKey name of dd api env var
const DDApiKey string = "DATADOG_API_KEY"

//DDAppKey name of dd app env var
const DDAppKey string = "DATADOG_APP_KEY"

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
	_, httpResponse, err := apiClient.AuthenticationApi.Validate(ctx).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when validating the DD auth keys: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", httpResponse)
		log.Fatal("\n--\nERROR Validating the API Call - check your env variables DATADOG_API_KEY and DATADOG_APP_KEY are set and are valid!") //clean this up
	}

	//fmt.Fprintf(os.Stdout, "API Call Validated\n\n")

	return ctx, apiClient

}
