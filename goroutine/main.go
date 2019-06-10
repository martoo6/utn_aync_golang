package main

//Gets data in parallel using goroutines

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

	avgHumidity := 0.0
	for _, h := range humidities {
		avgHumidity += h
	}
	avgHumidity /= float64(len(humidities))

	fmt.Println(avgHumidity)
}

func getLocationIds() []int {
	n := time.Now().UnixNano()
	//8 locations -> 8*500ms concurrently = 500ms + internal latency
	locations := []string{"london", "buenos aires", "rome", "tokyo", "new york", "madrid", "seoul", "moscow"}

	locationIds := make([]int, len(locations))
	var wgLocations sync.WaitGroup
	for i, location := range locations {
		wgLocations.Add(1)
		go func(index int, location string) {
			defer wgLocations.Done()
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

			locationIds[index] = locationsResp[0].Woeid
		}(i, location)
	}
	wgLocations.Wait()

	fmt.Printf("Locations all req latency: %v\n", time.Duration(time.Now().UnixNano()-n))
	return locationIds
}

func gethumidites(locationIds []int) []float64 {
	n := time.Now().UnixNano()

	var humiditiesMutex sync.Mutex
	var wgHumidities sync.WaitGroup
	var humidities []float64
	for _, locationId := range locationIds {
		wgHumidities.Add(1)
		go func(id int) {
			defer wgHumidities.Done()
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
				humiditiesMutex.Lock()
				humidities = append(humidities, w.Humidity)
				humiditiesMutex.Unlock()
			}
		}(locationId)
	}
	wgHumidities.Wait()

	fmt.Printf("Humidities all req latency: %v\n", time.Duration(time.Now().UnixNano()-n))
	return humidities
}
