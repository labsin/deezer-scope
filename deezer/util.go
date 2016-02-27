package main

import (
	//"log"
	"math"
	"strconv"
)

func PadWith0(number int, max int) string {
	var max_log int = int(math.Floor(math.Log10(float64(max))))
	number_log := 0
	if number > 0 {
		number_log = int(math.Floor(math.Log10(float64(number))))
	}
	var diff = max_log - number_log
	ret := strconv.Itoa(number)
	for diff > 0 {
		ret = "0" + ret
		diff--
	}
	return ret
}

func SecondsToString(seconds int) string {
	chunck := []int{
		60 * 60 * 24 * 365,
		60 * 60 * 24 * 30,
		60 * 60 * 24 * 7,
		60 * 60 * 24,
		60 * 60,
		60,
		1,
	}
	names := []string{
		`years`,
		`month`,
		`week`,
		`days`,
		`hours`,
		`minutes`,
		`seconds`,
	}
	var ret string
	for i, secondsInCount := range chunck {
		name := names[i]
		count := int(math.Floor(float64(seconds) / float64(secondsInCount)))
		if count > 0 {
			if ret != "" {
				ret += " "
			}
			ret += strconv.Itoa(count) + " " + name
			seconds -= count * secondsInCount
		}
	}
	if ret == "" {
		ret = "0 seconds"
	}
	return ret
}
