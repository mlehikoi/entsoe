package entsoe

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const endpoint = "https://web-api.tp.entsoe.eu/api"
const apiDateFormat = "200601020000"

type PricePoint struct {
	Time  time.Time
	Price float64
}

func GetPrices(token string, periodStart, periodEnd time.Time) ([]PricePoint, error) {
	url := fmt.Sprintf("%s?"+
		"securityToken=%s&"+
		"documentType=A44&"+
		"processType=A01&"+
		"in_Domain=10YFI-1--------U&"+
		"out_Domain=10YFI-1--------U&"+
		"periodStart=%s&"+
		"periodEnd=%s",
		endpoint, token, periodStart.Format(apiDateFormat), periodEnd.Format(apiDateFormat))

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Print the response body as a string
	// fmt.Println("Response Body:", string(body))

	var doc publicationMarketDocument
	err = xml.Unmarshal([]byte(body), &doc)
	if err != nil {
		fmt.Printf("Error unmarshaling XML: %v", err)
		return nil, err
	}

	var prices []PricePoint
	for _, ts := range doc.TimeSeries {
		t, err := time.Parse("2006-01-02T15:04Z", ts.Period.TimeInterval.Start)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		for _, point := range ts.Period.Points {
			price, err := strconv.ParseFloat(point.Amount, 64)
			if err != nil {
				fmt.Printf("Error converting string to float: %v", err)
				return nil, err
			}
			price /= 10.0
			price *= 1.24
			prices = append(prices, PricePoint{t, price})
			t = t.Add(time.Hour)
		}
	}
	return prices, nil
}

type publicationMarketDocument struct {
	XMLName    xml.Name     `xml:"Publication_MarketDocument"`
	TimeSeries []timeSeries `xml:"TimeSeries"`
}

type timeInterval struct {
	Start string `xml:"start"`
	End   string `xml:"end"`
}

type timeSeries struct {
	MRID             string `xml:"mRID"`
	BusinessType     string `xml:"businessType"`
	PriceMeasureUnit struct {
		Name string `xml:"name"`
	} `xml:"price_Measure_Unit"`
	Period struct {
		TimeInterval timeInterval `xml:"timeInterval"`
		Resolution   string       `xml:"resolution"`
		Points       []struct {
			Position int    `xml:"position"`
			Amount   string `xml:"price.amount"`
		} `xml:"Point"`
	} `xml:"Period"`
}
