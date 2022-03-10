package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func conversionTool(input, units string) (err error) {
	i, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return err
	}
	i = math.Round(i*10) / 10

	if strings.ToUpper(units) == "C" {
		c := (i - 32) * 5 / 9
		fmt.Printf("%.1fF = %.1fC\n", i, c)
	} else if strings.ToUpper(units) == "F" {
		c := (i * 9 / 5) + 32
		fmt.Printf("%.1fC = %.1fF\n", i, c)
	} else {
		return errors.New("invalid unit passed to conversion tool")
	}

	return err
}

func fakeForecastTool(input string) {
	fmt.Printf("The current temperature is %s degrees.\n", input)
}
