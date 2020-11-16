package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/brian-provenzano/dd-downtime/authentication"
)

//DefaultDowntimeMessage - Default downtime message if none is provided
const DefaultDowntimeMessage string = "Downtime scheduled by CEPSRE"

func main() {

	// define command line args for parsing
	parser := argparse.NewParser("commands", "Set Datadog downtimes for monitors via DD API.  Supports get, create, delete downtimes. "+
		"Requires DATADOG_APP_KEY and DATADOG_API_KEY environment variables.")

	//raw output top level option
	//rawFlag := parser.Flag("r", "raw", &argparse.Options{Required: false,
	//Help: "Limit output to raw json output from API response.  Useful if you wish to parse the response."})

	listCmd := parser.NewCommand("get", "List all the current Datadog downtimes")

	createCmd := parser.NewCommand("create", "Creates / Schedules the Datadog downtime")
	createCmdScope := createCmd.String("s", "scope",
		&argparse.Options{
			Required: true,
			Help:     "Existing datadog scope tag(s) - e.g. 'environment:prd1,service:voice-platform' (required)",
		})
	createCmdTime := createCmd.Int("t", "time",
		&argparse.Options{
			Required: true,
			Help:     "Downtime 'time' to set in minutes (required)",
		})
	createCmdMessage := createCmd.String("m", "message",
		&argparse.Options{
			Required: false,
			Help:     "Existing datadog scope tag(s) - e.g. 'environment:prd1,service:voice-platform' (required)",
		})

	updateCmd := parser.NewCommand("update", "Updates the Datadog downtime with the provided ID")
	updateCmdID := updateCmd.Int("i", "id",
		&argparse.Options{
			Required: false,
			Help:     "The ID of the Datadog downtime to update (required)",
		})

	deleteCmd := parser.NewCommand("delete", "Deletes the Datadog downtime with the provided ID")
	deleteCmdID := deleteCmd.Int("i", "id",
		&argparse.Options{
			Required: false,
			Help:     "The ID of the Datadog downtime to update (required)",
		})

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
			//log.Fatal("ERROR")
			return
		}

		//fmt.Fprintf(os.Stdout, "Response from `ListDowntimes`: %v\n", resp)

		//TODO - stub
		for i, s := range resp {
			fmt.Println(i, s.GetActive(), s.GetMessage())
		}

	case createCmd.Happened():
		//TODO - stub
		fmt.Println("This would run create based on other parameters (scope, message, time etc)")
		fmt.Fprintf(os.Stderr, "scope for create: %s\n", *createCmdScope)
		fmt.Fprintf(os.Stderr, "time in minutes for create: %d\n", *createCmdTime)
		if len(*createCmdMessage) > 0 {
			fmt.Fprintf(os.Stderr, "message to use for create: %s\n", *createCmdMessage)
		}

	case updateCmd.Happened():
		//TODO - stub
		fmt.Println("This would run update based on ID flag and other parameters (scope, message, endtime etc)")
		fmt.Fprintf(os.Stderr, "id for update: %d\n", *updateCmdID)

	case deleteCmd.Happened():
		//TODO - stub
		fmt.Println("This would run delete based on ID flag and other parameters (scope, message, endtime etc)")
		fmt.Fprintf(os.Stderr, "id for delete: %d\n", *deleteCmdID)
	}

}
