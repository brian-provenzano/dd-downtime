package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brian-provenzano/dd-downtime/authentication"
)

func main() {

	ctx, apiClient := authentication.Authenticate()

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
}
