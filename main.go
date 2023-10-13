package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Args[1]
	endpoint := "https://web-api.tp.entsoe.eu/api"
	// Construct the request URL for day-ahead prices in Finland
	// documentType A44 - Price document
	// processType A01 - Day ahead
	url := fmt.Sprintf("%s?securityToken=%s&documentType=A44&processType=A01&in_Domain=10YFI-1--------U&out_Domain=10YFI-1--------U&periodStart=202310090000&periodEnd=202310100000", endpoint, apiKey)

	// Make the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body as a string
	fmt.Println("Response Body:", string(body))
}
