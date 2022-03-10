package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// getInput reads user input from the command line and returns it (all lower case).
func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scanner.Scan()
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return strings.ToLower(scanner.Text())
}

// displayLogo reads the contents of the passed logo file and
// prints to the terminal.
func greeting() error {
	file, err := ioutil.ReadFile(logoFile)
	if err != nil {
		return err
	}
	fmt.Println(string(file))

	fmt.Println("Hello! Welcome to go-weather.")

	return nil
}

// executeMainMenu gets user input and calls the appropriate function.
func executeMainMenu() {
	fmt.Println("\nMAIN | Select an option:")
	for i, o := range mainMenu {
		fmt.Printf("%d) %s\n", i+1, o)
	}
Loop:
	for {
		input := getInput()
		switch input {
		case "1", mainMenu[0]:
			executeWeatherMenu()
			break Loop
		case "2", mainMenu[1]:
			executeConversion()
			break Loop
		case "3", mainMenu[2]:
			executeGenerator()
			break Loop
		case "4", strings.ToLower(mainMenu[3]):
			quitProgram()
			break Loop
		default:
			fmt.Println("Invalid input. Enter option number or name.")
		}
	}
}

// executeWeatherMenu gets user input and executes the correct weather request.
func executeWeatherMenu() {
	fmt.Println("\nWEATHER | Select an option:")
	for i, o := range weatherMenu {
		fmt.Printf("%d) %s\n", i+1, o)
	}
Loop:
	for {
		input := getInput()
		switch input {
		case "1", weatherMenu[0]:
			executeWeatherRequest("1")
			break Loop
		case "2", weatherMenu[1]:
			executeWeatherRequest("2")
			break Loop
		case "3", weatherMenu[2]:
			executeWeatherRequest("3")
			break Loop
		case "4", strings.ToLower(weatherMenu[3]):
			break Loop
		default:
			fmt.Println("Invalid input. Enter option number or name.")
		}
	}
}

// executeConversion converts values between Celsius and Fahrenheit and outputs the result.
func executeConversion() {
	fmt.Println("Enter units to convert to (C or F):")
	units := strings.ToUpper(getInput())
	if units != "C" && units != "F" {
		fmt.Println("Invalid units. Must enter 'C' for Celsius or 'F' for Fahrenheit.")
		time.Sleep(2 * time.Second)
		return
	}
	fmt.Println("Enter value to convert: ")
	value := getInput()
	err := conversionTool(value, units)
	if err != nil {
		fmt.Println("Invalid input. Must enter numerical value.")
	}

	time.Sleep(2 * time.Second)
}

// executeGenerator outputs a fake forecast as a string containing the value received
// from the user input.
func executeGenerator() {
	fmt.Println("Enter your ideal temperature (ex. 68):")
	input := getInput()
	_, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Must enter numerical value.")
	} else {
		fakeForecastTool(input)
	}

	time.Sleep(2 * time.Second)
}

// executeWeatherRequest returns the results of one call to the OpenWeatherMap API for a location
// for the period of time determined by the passed option.
func executeWeatherRequest(option string) {
	fmt.Println("Enter the partial/full address of where you want to get the weather:")
	input := getInput()
	place, err := getCoordForPlace(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch option {
	case "1":
		weather, err := getWeather(place, UnitsImperial, "current")
		if err != nil {
			fmt.Println(err)
			return
		}
		weather.Current.Output(UnitsImperial)
	case "2":
		weather, err := getWeather(place, UnitsImperial, "daily")
		if err != nil {
			fmt.Println(err)
			return
		}
		weather.outputWeeklyForecast(UnitsImperial)
	case "3":
		weather, err := getWeather(place, UnitsImperial, "daily")
		if err != nil {
			fmt.Println(err)
			return
		}
		weather.outputAvgWeather(UnitsImperial)
	default:
		fmt.Println("Invalid option passed to executeWeatherRequest.")
	}

	time.Sleep(2 * time.Second)
}
