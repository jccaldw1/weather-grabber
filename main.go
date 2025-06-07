package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DailyUnits struct {
	Time               string `json:"time"`
	Temperature_2m_Max string `json:"temperature_2m_max"`
	Temperature_2m_Min string `json:"temperature_2m_min"`
}

type DailyResults struct {
	Time               []string  `json:"time"`
	Temperature_2m_Max []float32 `json:"temperature_2m_max"`
	Temperature_2m_Min []float32 `json:"temperature_2m_min"`
}

type WeatherResponse struct {
	Latitude             float32      `json:"latitude"`
	Longitude            float32      `json:"longitude"`
	GenerationTime       float64      `json:"generationtime_ms"`
	UtcOffsetSeconds     int          `json:"utc_offset_seconds"`
	TimeZone             string       `json:"timezone"`
	TimeZoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float32      `json:"elevation"`
	DailyUnits           DailyUnits   `json:"daily_units"`
	Daily                DailyResults `json:"daily"`
}

type DatabaseObjectToPost struct {
	Date      string
	High      float32
	Low       float32
	DaysAhead int
}

func main() {
	fmt.Println("Hello world")

	url := "https://api.open-meteo.com/v1/forecast?latitude=33.4148&longitude=-111.9093&daily=temperature_2m_max,temperature_2m_min&hourly=temperature_2m&timezone=auto&forecast_days=14&temperature_unit=fahrenheit"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching reqeust: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return
	}

	var weatherResponse WeatherResponse

	err = json.Unmarshal([]byte(body), &weatherResponse)

	if err != nil {
		fmt.Println("Error unmarshalling response")
		return
	}

	fmt.Printf("Response: %+v\n", weatherResponse)

	var databaseObjectList []DatabaseObjectToPost

	for i := 0; i < len(weatherResponse.Daily.Time); i++ {
		var databaseObject DatabaseObjectToPost
		databaseObject.Date = weatherResponse.Daily.Time[i]
		databaseObject.DaysAhead = i
		databaseObject.High = weatherResponse.Daily.Temperature_2m_Max[i]
		databaseObject.Low = weatherResponse.Daily.Temperature_2m_Min[i]

		databaseObjectList = append(databaseObjectList, databaseObject)
	}

	for i := 0; i < len(databaseObjectList); i++ {
		fmt.Printf("Database object: %+v\n", databaseObjectList[i])
	}
}
