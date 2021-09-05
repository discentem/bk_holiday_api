package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/golang/gddo/httputil/header"
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

func mapOfHolidays(holidays *[]Holiday) map[string]Holiday {
	hm := map[string]Holiday{}
	for _, holiday := range *holidays {
		hm[*holiday.Date] = holiday
	}
	return hm
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
		logger.Error(string(byt))
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

func IsHolidayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryCode := vars["countryCode"]
	date := vars["date"]
	year := strings.Split(date, "-")
	logger.Infof("year: %s", year)
	holidays, err := getPublicHolidays(year[0], countryCode)
	if err != nil {
		logger.Error(err)
	}
	hm := mapOfHolidays(holidays)
	if val, ok := hm[date]; ok {
		byt, err := json.Marshal(val)
		if err != nil {
			logger.Error(err)
		}
		msg := fmt.Sprintf("Yes %s is a holiday in %s\n", date, countryCode)
		w.Write([]byte(msg))
		w.Write(byt)
		return
	}
	msg := fmt.Sprintf("no holidays on %s in %s", date, countryCode)
	w.Write([]byte(msg))
}

var (
	serverURL = "localhost:8080"
)

func AreTheseDates(dates []string, countryCode string) ([]byte, error) {
	data := []byte{}
	for _, date := range dates {
		url := fmt.Sprintf("http://%s/isHoliday/%s/%s", serverURL, date, countryCode)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		byt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		byt = append(byt, []byte("\n")...)
		data = append(data, byt...)
	}
	return data, nil
}

func AreTheseHolidaysHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryCode := vars["countryCode"]
	dates := strings.Split(vars["dates"], ",")
	if len(dates) == 0 {
		w.Write([]byte("no dates provided"))
		logger.Error("no dates provided")
	}
	byt, err := AreTheseDates(dates, countryCode)
	if err != nil {
		http.Error(w, err.Error(), 1)
	}
	w.Write(byt)
}

func checkHeaderIsJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}
}

func AreTheseHolidaysJSONHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryCode := vars["countryCode"]
	checkHeaderIsJSON(w, r)
	type DatesList struct {
		Dates []string `json:"dates"`
	}
	var dates DatesList
	err := json.NewDecoder(r.Body).Decode(&dates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info(dates)
	byt, err := AreTheseDates(dates.Dates, countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(byt)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/holidays/{year}/{countryCode}", HolidaysHandler)
	r.HandleFunc("/isHoliday/{date}/{countryCode}", IsHolidayHandler)
	// dates expected to be comma separated
	r.HandleFunc("/areTheseHolidays/{countryCode}/{dates}", AreTheseHolidaysHandler)
	r.HandleFunc("/areTheseHolidaysJSON/{countryCode}", AreTheseHolidaysJSONHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logger.Fatal(http.ListenAndServe(serverURL, r))
	}()
	//time.Sleep(time.Second * 5)

	var jsonStr = []byte(`{"dates":["2021-07-05"]}`)
	uri := fmt.Sprintf("http://%s/areTheseHolidaysJSON/US", serverURL)
	req, err := http.NewRequest("GET", uri, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	wg.Wait()
}
