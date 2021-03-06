package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ===========================================
// The majority of this code has been authored
// by David Cako. All additions and edits have
// been marked with the appropriate comment.
// ===========================================

// OpenWeatherCondition stores the values for the 'Weather' object in all responses.
type OpenWeatherCondition struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

// OpenWeatherResponseCurrent stores the weather values for the Current response.
type OpenWeatherResponseCurrent struct {
	// unix stamp - number of ms that have passed since 1 Dec 1969
	Dt         int64
	Sunrise    int64
	Sunset     int64
	Temp       float32
	Feels_like float32
	Pressure   int
	Humidity   int
	Dew_point  float32
	Uvi        float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       struct {
		_1hr float32 `json:"1hr"`
	}
	Snow struct {
		_1hr float32 `json:"1hr"`
	}
}

// (OpenWeatherResponseCurrent) Output returns a string containing current weather values.
func (o OpenWeatherResponseCurrent) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	return fmt.Sprintf("Temp: %.1f%s | Humidity: %d%% | %s\n", o.Temp, unitAbbr, o.Humidity, o.Weather[0].Description)
}

// OpenWeatherResponseHourly stores the weather values for the Hourly responses.
type OpenWeatherResponseHourly struct {
	Dt         int64
	Temp       float32
	Feels_like float32
	Pressure   int
	Humidity   int
	Dew_point  float32
	//Uvi        float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       struct {
		_1hr float32 `json:"1hr"`
	}
	Snow struct {
		_1hr float32 `json:"1hr"`
	}
}

// (OpenWeatherResponseHourly) Output returns a string containing hourly weather values.
func (o OpenWeatherResponseHourly) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	t := time.Unix(o.Dt, 0)

	return fmt.Sprintf("%-9s %2d/%2d %02d:00\t%5.2f%s | Humidity: %d%% | %s\n", t.Weekday().String(),
		t.Month(), t.Day(), t.Hour(), o.Temp, unitAbbr, o.Humidity, o.Weather[0].Description)
}

// OpenWeatherResponseDaily stores the weather values for the Daily responses.
type OpenWeatherResponseDaily struct {
	Dt      int64
	Sunrise int64
	Sunset  int64
	Temp    struct {
		Day   float32
		Min   float32
		Max   float32
		Night float32
		Eve   float32
		Morn  float32
	}
	Feels_like struct {
		Day   float32
		Night float32
		Eve   float32
		Morn  float32
	}
	Pressure   int
	Humidity   int
	Dew_point  float32
	Uvi        float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       float32
	Snow       float32
}

// (OpenWeatherResponseDaily) Output returns a string containing daily weather values.
func (o OpenWeatherResponseDaily) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	t := time.Unix(o.Dt, 0)

	return fmt.Sprintf("%-9s %2d/%2d\tHigh: %5.2f%s Low: %5.2f%s | Humidity: %d%% | %s\n", t.Weekday().String(),
		t.Month(), t.Day(), o.Temp.Max, unitAbbr, o.Temp.Min, unitAbbr, o.Humidity, o.Weather[0].Description)
}

// OpenWeatherResponseOneCall stores all responses received from OpenWeatherMap API requests.
type OpenWeatherResponseOneCall struct {
	Current *OpenWeatherResponseCurrent
	Hourly  *[]OpenWeatherResponseHourly
	Daily   *[]OpenWeatherResponseDaily
}

// getWeather sends a single request to the OpenWeatherMap API for the weather
// at the passed coordinates. Responses are stored in an OpenWeatherResponseOneCall
// variable.
func GetWeather(c Coordinates, units, period string) (weather OpenWeatherResponseOneCall, err error) {
	exclude := []string{WeatherPeriodMinutely}

	if period != WeatherPeriodCurrent {
		exclude = append(exclude, WeatherPeriodCurrent)
	}
	if period != WeatherPeriodHourly {
		exclude = append(exclude, WeatherPeriodHourly)
	}
	if period != WeatherPeriodDaily {
		exclude = append(exclude, WeatherPeriodDaily)
	}

	excString := strings.Join(exclude, ",")

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?appid=%s&lat=%g&lon=%g&exclude=%s&units=%s",
		keys.OpenWeatherMapApiKey, c.Lat, c.Lng, excString, units)

	r, err := httpClient.Get(url)
	if err != nil {
		return weather, err
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		return weather, fmt.Errorf("OpenWeatherRequest failed: %s", r.Status)
	}

	err = json.NewDecoder(r.Body).Decode(&weather)

	return weather, err
}

// ADDITION: added AverageWeather

// AverageWeather stores the average weather values for a 7-day week.
type AverageWeather struct {
	MinTemp  float32
	MaxTemp  float32
	Pressure int
	Humidity int
	Uvi      float32
}

// ADDITION: added PrintAvgWeather

// PrintAvgWeather calculates the average weather values for a 7-day week and prints
// the output.
func (o OpenWeatherResponseOneCall) PrintAvgWeather(units string) {
	var min, max, uvi float32
	var prs, hum int
	for i, day := range *o.Daily {
		if i == 8 {
			break
		} else if i != 0 {
			min += day.Temp.Min
			max += day.Temp.Max
			uvi += day.Uvi
			prs += day.Pressure
			hum += day.Humidity
		}
	}
	avg := AverageWeather{min / 7, max / 7, prs / 7, hum / 7, uvi / 7}

	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	fmt.Printf("AVG WEEKLY WEATHER\nMin: %.1f%s Max: %.1f%s | Pressure: %d mb | Humidity: %d%% | UVI: %.1f\n",
		avg.MinTemp, unitAbbr, avg.MaxTemp, unitAbbr, avg.Pressure, avg.Humidity, avg.Uvi)
}

// ADDITION: added PrintWeeklyForecast

// PrintWeeklyForecast generates a 7-day weather forecast and prints the output.
func (o OpenWeatherResponseOneCall) PrintWeeklyForecast(units string) {
	fmt.Println("WEEKLY FORECAST")

	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	var t time.Time
	for i, day := range *o.Daily {
		if i == 8 {
			break
		} else if i != 0 {
			t = time.Unix(day.Dt, 0)
			fmt.Printf("%-9s %2d/%2d  -  Low: %.1f%s High: %.1f%s | Humidity: %d%% | %s\n", t.Weekday().String(),
				t.Month(), t.Day(), day.Temp.Min, unitAbbr, day.Temp.Max, unitAbbr, day.Humidity, day.Weather[0].Description)
		}
	}
}

// ADDITION: added PrintCurrent

// (OpenWeatherResponseCurrent) Output prints the current weather values.
func (o OpenWeatherResponseCurrent) PrintCurrent(units string) {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	fmt.Printf("CURRENT WEATHER\nTemp: %.1f%s | Humidity: %d%% | %s\n", o.Temp, unitAbbr, o.Humidity, o.Weather[0].Description)
}
