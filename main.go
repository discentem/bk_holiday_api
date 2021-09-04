package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/logger"
	"github.com/kr/pretty"
)

type Holiday struct {
	Date        *string   `json:"date"`
	LocalName   *string   `json:"localName"`
	Name        *string   `json:"name"`
	CountryCode *string   `json:"countryCode"`
	Fixed       bool      `json:"fixed"`
	Global      bool      `json:"global"`
	Counties    *[]string `json:"counties"`
	LaunchYear  *int      `json:"launchYear"`
	Types       *[]string `json:"types"`
}

func getPublicHolidays() (*[]Holiday, error) {
	resp, err := http.Get("https://date.nager.at/api/v3/publicholidays/2021/AT")
	if err != nil {
		return nil, err
	}
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp.Status was: %d", resp.StatusCode)
	}
	var holidays *[]Holiday
	err = json.Unmarshal(byt, &holidays)
	if err != nil {
		return nil, err
	}

	return holidays, nil

}

func init() {
	logger.Init("bk_holiday_api", true, false, io.Discard)
}

func main() {
	holidays, err := getPublicHolidays()
	if err != nil {
		logger.Info(pretty.Print(holidays))
		logger.Fatal(err)
	}
	logger.Info(pretty.Print(holidays))
}
