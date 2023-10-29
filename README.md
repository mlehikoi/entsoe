# entsoe

entsoe is a Go API designed for retrieving day-ahead electricity prices from ENTSO-E (European association for the cooperation of transmission system operators for electricity).
This API specifically provides prices for the Finnish market.

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

The prices are provided in cents per kilowatt-hour (kWh) and include Finnish Value Added Tax (VAT).

## Obtaining the token

To use the API, you need an authentication token.
Here are the steps to obtain it:

1. **Register on ENTSO-E:** Go to https://transparency.entsoe.eu/ and complete the registration process.
2. **Request API Access:** After registering, send an email to request access to the RESTful API.
3. **For More Information:** Detailed information about authentication and authorization can be found at https://transparency.entsoe.eu/content/static_content/Static%20content/web%20api/Guide.html#_authentication_and_authorisation.
