# entsoe

entsoe is a Go API for retrivieving day-ahead prices from ENTSO-E (European association for the cooperation of transmission system operators (TSOs) for electricity).
It provides the prices for the Finnish market.

## Usage example

    package main

    import (
        "fmt"
        "log"
        "os"
        "time"

        "github.com/mlehikoi/entsoe"
    )

    func main() {
        now := time.Now()
        beginPeriod := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-24 * time.Hour)
        endPeriod := now.Add(24 * time.Hour)

        token := os.Getenv("ENTSOE_TOKEN")
        prices, _ := entsoe.GetPrices(token, beginPeriod, endPeriod)

        localTimeZone, err := time.LoadLocation("Local")
        if err != nil {
            log.Fatalln("Error getting location:", err)
            return
        }
        n := 0
        sum := 0.0
        for _, pp := range prices {
            local := pp.Time.In(localTimeZone)
            if local.Month() == now.Month() {
                sum += pp.Price
                n += 1
            }
        }
        fmt.Printf("Monthly average price: %.2f\n", sum/float64(n))
    }

## Notes

The prices are in cents per kWh and include the Finnish VAT.

## Getting the token

In order to use the API, you need a token.
Register to entsoe: https://transparency.entsoe.eu/
After having registered, send an email to get access to the RESTful API.
For more information, see: https://transparency.entsoe.eu/content/static_content/Static%20content/web%20api/Guide.html#_authentication_and_authorisation
