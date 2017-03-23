/*

Simple tool to fetch data from the FatSecret API and dump as JSON

The following environment variables must be set in order to use the
FatSecret API...

1) "FATSECRET_CONSUMER_KEY"
2) "FATSECRET_SHARED_SECRET"

*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fitzone/fatsecret"
	"strings"
)

var (
	consumerKey  = os.Getenv("FATSECRET_CONSUMER_KEY")
	sharedSecret = os.Getenv("FATSECRET_SHARED_SECRET")
)

func main() {
	// define and parse the cli flags
	methodPtr := flag.String("method", "food_categories.get", "FatSecret API method")
	paramsPtr := flag.String("params", "", "Comma-separated FatSecret API parameters")
	flag.Parse()

	// build the api parameters map
	params := map[string]string{}
	if len(*paramsPtr) > 0 {
		for _, p := range strings.Split(*paramsPtr, ",") {
			tokens := strings.Split(p, "=")
			if len(tokens) != 2 {
				fmt.Printf("\nInvalid API name/value parameter given ('%v'). Must use format 'name=value'.\n\n", p)
				os.Exit(1)
			}
			params[tokens[0]] = tokens[1]
		}
	}

	// create a fatsecret client
	client, err := fatsecret.NewClient(consumerKey, sharedSecret)
	if err != nil {
		panic(err)
	}

	// invoke the low-level fatsecret api
	body, err := client.InvokeAPI(*methodPtr, params)
	if err != nil {
		fmt.Printf("Error invoking the '%s' FatSecret API: err = '%v'", *methodPtr, err)
	} else {
		// dump the api json response
		fmt.Println(string(body))
	}
}
