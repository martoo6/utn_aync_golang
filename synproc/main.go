package main

//Gets data sequentially

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LocationsResp []LocationResp

type LocationResp struct {
	Woeid int `json:"woeid"`
}

type LocationWeather struct {
	ConsolidatedWeather []struct {
		AirPressure float64 `json:"air_pressure"`
		Humidity    float64 `json:"humidity"`
		Visibility  float64 `json:"visibility"`
	} `json:"consolidated_weather"`
}

var apiUrl = "http://localhost:8080/"

func main() {
	locationIds := getLocationIds()
	humidities := gethumidites(locationIds)

	var avgHumidity float64
	for _, h := range humidities {
		avgHumidity += h
	}
	avgHumidity /= float64(len(humidities))

	fmt.Println(avgHumidity)
}

func getLocationIds() []int {
	n := time.Now().UnixNano()
	//8 locations -> 8*500ms = 4 seconds + internal latency
	locations := []string{"london", "buenos aires", "rome", "tokyo", "new york", "madrid", "seoul", "moscow"}

	locationIds := make([]int, len(locations))
	for i, location := range locations {
		//Make request
		time.Sleep(500 * time.Millisecond)
		resp, _ := http.Get(apiUrl + "location/search/?query=" + url.QueryEscape(location))

		//Read body bytes
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))

		//Unmarshall
		var locationsResp LocationsResp
		json.Unmarshal(body, &locationsResp)

		locationIds[i] = locationsResp[0].Woeid
	}
	fmt.Printf("Locations all req latency: %v\n", time.Duration(time.Now().UnixNano()-n))

	return locationIds
}

func gethumidites(locationIds []int) []float64 {
	n := time.Now().UnixNano()
	var humidities []float64
	for _, id := range locationIds {
		//Make request
		resp, _ := http.Get(apiUrl + "location/" + strconv.Itoa(id) + "/")

		//Read body bytes
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))

		//Unmarshall
		var locationWeather LocationWeather
		json.Unmarshal(body, &locationWeather)
		for _, w := range locationWeather.ConsolidatedWeather {
			humidities = append(humidities, w.Humidity)
		}
	}
	fmt.Printf("Humidities all req latency: %v\n", time.Duration(time.Now().UnixNano()-n))

	return humidities
}
