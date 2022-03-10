package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Constants used throughout main package.
const (
	WeatherPeriodCurrent  = "current"
	WeatherPeriodMinutely = "minutely"
	WeatherPeriodHourly   = "hourly"
	WeatherPeriodDaily    = "daily"
	UnitsImperial         = "imperial"
	UnitsMetric           = "metric"
	keysFile              = ".apiConfig"
	logoFile              = "logo.txt"
	menuFile              = ".menuConfig"
)

// apiKeys stores the API keys used in the main package.
type apiKeys struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
	GoogleApiKey         string `json:"GoogleApiKey"`
}

// keys stores a reference to the API keys.
var keys apiKeys

// httpClient is used to make requests to APIs.
var httpClient http.Client

// mainMenu stores the options for the first menu displayed.
var mainMenu = []string{"weather", "conversion tool", "fake forecast generator", "QUIT"}

// weatherMenu stores the options for the 'weather' sub-menu.
var weatherMenu = []string{"current weather", "7-day forecast", "weekly averages", "BACK"}

func main() {
	// load API keys from JSON
	err := loadKeys()
	if err != nil {
		panic(err)
	}

	// start HTTP client
	httpClient = http.Client{
		Timeout: time.Second * 10,
	}

	// print greeting from TXT
	err = greeting()
	if err != nil {
		panic(err)
	}

	// start program loop
	loop()
}

func loop() {
	for {
		// display main menu
		executeMainMenu()
	}
}

func quitProgram() {
	fmt.Println("Goodbye..")
	fmt.Println()
	os.Exit(0)
}

// loadKeys parses the API keys in the .apiConfig file and assigns
// them to the apiKeys variable.
func loadKeys() error {
	bytes, err := ioutil.ReadFile(keysFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return err
	}

	return nil
}
