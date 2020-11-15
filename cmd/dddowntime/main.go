package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"

	"github.com/brian-provenzano/dd-downtime/authentication"
)

func main() {

	// define command line args for parsing
	parser := argparse.NewParser("commands", "Simple example of argparse commands")

	listCmd := parser.NewCommand("list", "List all the current Datadog downtimes")
	updateCmd := parser.NewCommand("update", "Updates the Datadog downtime with the provided ID")

	// Boilerplate: Parse command line arguments and in case of any error print error and help information
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	//authenticate to DD API
	ctx, apiClient := authentication.Authenticate()

	switch {

	case listCmd.Happened():
		resp, r, err := apiClient.DowntimesApi.ListDowntimes(ctx).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ListDowntimes`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			log.Fatal("ERROR")
		}

		fmt.Fprintf(os.Stdout, "Response from `ListDowntimes`: %v\n", resp)

		for i, s := range resp {
			fmt.Println(i, s.GetActive(), s.GetMessage())
		}

	case updateCmd.Happened():
		fmt.Println("This would run update based on ID flag and other parameters (scope, message, endtime etc)")

	}

}
