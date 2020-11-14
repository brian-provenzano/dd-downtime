package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brian-provenzano/dd-downtime/authentication"
)

func main() {

	ctx, apiClient := authentication.Authenticate()

	dashboardID := "<DASHBOARD_ID>"

	resp, r, err := apiClient.DashboardsApi.GetDashboard(ctx, dashboardID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DashboardsApi.GetDashboard`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		log.Fatal("ERROR")
	}

	fmt.Fprintf(os.Stdout, "Response from `DashboardsApi.GetDashboard`: %v\n", resp)

}
