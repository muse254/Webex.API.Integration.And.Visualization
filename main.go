package main

import "Webex.API.Integration.And.Visualization/api"

func main() {
	// Start the server and close if error occurs
	if err := api.WebexApplicationServer(); err != nil {
		panic(err)
	}
}
