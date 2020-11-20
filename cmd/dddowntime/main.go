package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/akamensky/argparse"
	"github.com/davecgh/go-spew/spew"

	"github.com/brian-provenzano/dd-downtime/authentication"
)

//DefaultDowntimeMessage - Default downtime message if none is provided
const DefaultDowntimeMessage string = "Downtime scheduled by CEPSRE"

func main() {

	// define command line args for parsing
	parser := argparse.NewParser("commands", "Get, Create, Update, Delete Datadog downtimes via DD API. "+
		"Requires DATADOG_APP_KEY and DATADOG_API_KEY environment variables to be set for authentication to the DD API.")

	//raw output top level option - TODO
	//rawFlag := parser.Flag("r", "raw", &argparse.Options{Required: false,
	//Help: "Limit output to raw json output from API response.  Useful if you wish to parse the response."})

	//debug option
	debugFlag := parser.Flag("d", "debug", &argparse.Options{Required: false,
		Help: "Show debug information for the action"})

	getCmd := parser.NewCommand("get", "Get/List Datadog downtime that corresponds to the provided ID")
	getCmdID := getCmd.Int("i", "id",
		&argparse.Options{
			Required: true,
			Help:     "The ID for the Datadog downtime to retrieve (required) - e.g. 123",
		})

	createCmd := parser.NewCommand("create", "Creates/Schedules the Datadog downtime")
	createCmdScope := createCmd.String("s", "scope",
		&argparse.Options{
			Required: true,
			Help:     "Existing datadog scope tag(s) as string - e.g. 'environment:prd1,service:voice-platform' (required)",
		})
	createCmdTime := createCmd.String("t", "time",
		&argparse.Options{
			Required: true,
			Help:     "Downtime 'time' to set in minutes as string (required) - e.g. 60m, 1h, 320s etc. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.",
		})
	createCmdMessage := createCmd.String("m", "message",
		&argparse.Options{
			Required: false,
			Help:     "Existing datadog scope tag(s) - e.g. 'environment:prd1,service:voice-platform' (required)",
		})

	updateCmd := parser.NewCommand("update", "Updates the Datadog downtime that corresponds to the provided ID")
	updateCmdID := updateCmd.Int("i", "id",
		&argparse.Options{
			Required: true,
			Help:     "The ID of the Datadog downtime to update (required) - e.g. 123",
		})

	updateCmdScope := updateCmd.String("s", "scope",
		&argparse.Options{
			Required: false,
			Help:     "Existing datadog scope tag(s) as string - e.g. 'environment:prd1,service:voice-platform'",
		})
	updateCmdTime := updateCmd.String("t", "time",
		&argparse.Options{
			Required: false,
			Help:     "Downtime 'time' to set in minutes as string (required) - e.g. 60m, 1h, 320s etc. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.",
		})
	updateCmdMessage := updateCmd.String("m", "message",
		&argparse.Options{
			Required: false,
			Help:     "Existing datadog scope tag(s) - e.g. 'environment:prd1,service:voice-platform'",
		})

	deleteCmd := parser.NewCommand("delete", "Deletes the Datadog downtime that corresponds to the provided ID")
	deleteCmdID := deleteCmd.Int("i", "id",
		&argparse.Options{
			Required: true,
			Help:     "The ID of the Datadog downtime to delete (required) - e.g. 123",
		})

	// Parse command line arguments and in case of any error print error and help information
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	//authenticate to DD API
	ctx, apiClient := authentication.Authenticate()

	switch {
	case getCmd.Happened():

		downtimeId := int64(*getCmdID)
		downtime, resp, err := apiClient.DowntimesApi.GetDowntime(ctx, downtimeId).Execute()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `GetDowntimes`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
			return
		}

		if *debugFlag {
			spew.Dump(downtime)
		}

		fmt.Fprintf(os.Stdout, "Downtime Details\n ID: %v \n Active?: %v \n Message: %v \n Scope: %v \n StartTime: %v \n EndTime: %v \n HasEnd?: %v \n",
			downtime.GetId(), downtime.GetActive(), downtime.GetMessage(), downtime.GetScope(), time.Unix(downtime.GetStart(), 0), time.Unix(downtime.GetEnd(), 0), downtime.HasEnd())
	case createCmd.Happened():

		body := *datadog.NewDowntime()
		//TODO - check format bc this split will fail
		body.SetScope(strings.Split(*createCmdScope, ","))

		if len(*createCmdMessage) > 0 {
			body.SetMessage(*createCmdMessage)
		} else {
			body.SetMessage(DefaultDowntimeMessage)
		}

		duration, err := time.ParseDuration(*createCmdTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid format for Time parameter: %v\nValid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.", err)
			return
		}
		downtimeEnd := time.Now().Add(duration)
		body.SetEnd(downtimeEnd.Unix())

		if *debugFlag {
			spew.Dump(body)
		}

		downtime, resp, err := apiClient.DowntimesApi.CreateDowntime(ctx).Body(body).Execute()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `CreateDowntime`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
			return
		}

		fmt.Fprintf(os.Stdout, "Downtime Created!\n ID: %v \n Active?: %v \n Message: %v \n Scope: %v \n StartTime: %v \n EndTime: %v \n HasEnd?: %v \n",
			downtime.GetId(), downtime.GetActive(), downtime.GetMessage(), downtime.GetScope(), time.Unix(downtime.GetStart(), 0), time.Unix(downtime.GetEnd(), 0), downtime.HasEnd())

	case updateCmd.Happened():

		downtimeId := int64(*updateCmdID)

		//TODO can prob reflect through the argsparse Command struct to this to clean up the optionCount stuff
		//there are no easy accessor methods ...
		//spew.Dump(updateCmd)

		body := *datadog.NewDowntime()
		optionCount := 0

		if len(*updateCmdScope) > 0 {
			//TODO check if string format not null/empty ; correct format
			body.SetScope(strings.Split(*updateCmdScope, ","))
			optionCount++
		}

		if len(*updateCmdTime) > 0 {
			duration, err := time.ParseDuration(*updateCmdTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid format for Time parameter: %v\nValid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.", err)
				return
			}
			downtimeEnd := time.Now().Add(duration)
			body.SetEnd(downtimeEnd.Unix())
			optionCount++
		}

		if len(*updateCmdMessage) > 0 {
			body.SetMessage(*updateCmdMessage)
			optionCount++
		}

		if *debugFlag {
			spew.Dump(body)
		}
		if optionCount == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: You must provide at least scope, time or message to update the downtime!")
			return
		}

		downtime, resp, err := apiClient.DowntimesApi.UpdateDowntime(ctx, downtimeId).Body(body).Execute()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `UpdateDowntime`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
			return
		}

		fmt.Fprintf(os.Stdout, "Downtime Updated!\n ID: %v \n Active?: %v \n Message: %v \n Scope: %v \n StartTime: %v \n EndTime: %v \n HasEnd?: %v \n",
			downtime.GetId(), downtime.GetActive(), downtime.GetMessage(), downtime.GetScope(), time.Unix(downtime.GetStart(), 0), time.Unix(downtime.GetEnd(), 0), downtime.HasEnd())

	case deleteCmd.Happened():

		downtimeId := int64(*deleteCmdID)

		//TODO - GET STUB :let's call 'get' first...display the details for downtime we are going to delete before we delete it

		resp, err := apiClient.DowntimesApi.CancelDowntime(ctx, downtimeId).Execute()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `DowntimesApi.CancelDowntime``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
		}

		fmt.Fprintf(os.Stderr, "Downtime with ID [ %d ] deleted successfully! \n", *deleteCmdID)
	}
}
