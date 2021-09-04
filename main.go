package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/logger"
	"github.com/gorilla/mux"
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

func getPublicHolidays(year string, countryCode string) (*[]Holiday, error) {
	url := fmt.Sprintf(
		"https://date.nager.at/api/v3/publicholidays/%s/%s",
		year,
		countryCode)
	resp, err := http.Get(url)
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

func HolidaysHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]
	countryCode := vars["countryCode"]
	holidays, err := getPublicHolidays(year, countryCode)
	if err != nil {
		logger.Error(err)
	}
	byt, err := json.Marshal(holidays)
	if err != nil {
		logger.Error(err)
	}
	w.Write(byt)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/holidays/{year}/{countryCode}", HolidaysHandler)
	logger.Fatal(http.ListenAndServe("localhost:8080", r))

}
