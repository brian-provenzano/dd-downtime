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

	//fmt.Fprintf(os.Stderr, "apikey: %v\n", os.Getenv(DDApiKey))
	//fmt.Fprintf(os.Stderr, "appkey: %v\n", os.Getenv(DDAppKey))

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
		log.Fatal("ERROR") //clean this up
	}

	fmt.Fprintf(os.Stdout, "API Call Validated\n\n")

	return ctx, apiClient

}
