package weathers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_weather/src/db"
	"log"
	"net/http"
	"strings"

	"go_weather/src/service/weather"

	"github.com/go-chi/chi"
)

const openWeatherAPIKey = "2e474ce4f674a14ecb86151dca5a4fe2"

// WController with dependencies
type WController struct {
	DB                    db.DatabaseInterf
	GetWeather            func(string, weather.OpenWeatherAPIData, *sql.DB) weather.OpenWeatherAPIData
	GetQueryCount         func(*sql.DB) (int, error)
	GetHighTemps          func(*sql.DB) []weather.TemperatureRow
	GetLowTemps           func(*sql.DB) []weather.TemperatureRow
	GetAverageTempByMonth func(int, *sql.DB) (weather.TemperatureRow, error)
	GetDaysByWeatherType  func(*sql.DB) []weather.MainWeatherRow
	GetOpenWeatherAPIData func(url string) (resp *http.Response, err error)
}

// GetWeatherHandler calls getWeather service to retrieve weather data for provided city
func (wc *WController) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	cityName := chi.URLParam(r, "city")
	var openWeatherAPIData weather.OpenWeatherAPIData
	weatherResponse, err := wc.GetOpenWeatherAPIData("http://api.openweathermap.org/data/2.5/weather?q=" + cityName + "&&units=imperial&appid=" + openWeatherAPIKey)
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer weatherResponse.Body.Close()
	err = json.NewDecoder(weatherResponse.Body).Decode(&openWeatherAPIData)
	if err != nil {
		println("error decoding json", err)
	}

	response := wc.GetWeather(cityName, openWeatherAPIData, wc.DB.OpenConnection())

	resp, err := json.Marshal(response)
	if err != nil {
		println("error converting ", err)
	}

	fmt.Println("resp ", resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetQueryCountHandler calls service to get total number of queries entered
func (wc *WController) GetQueryCountHandler(w http.ResponseWriter, r *http.Request) {
	queryCount, err := wc.GetQueryCount(wc.DB.OpenConnection())
	if err != nil {
		fmt.Println("error", err)
	}

	responseStruct := struct {
		QueryCount int `json:"queryCount"`
	}{
		queryCount,
	}
	response, err := json.Marshal(responseStruct)
	if err != nil {
		fmt.Println("err:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetHighTempsPerMonthHandler returns temperature highs for each month that has been querie
func (wc *WController) GetHighTempsPerMonthHandler(w http.ResponseWriter, r *http.Request) {
	highTemps := wc.GetHighTemps(wc.DB.OpenConnection())

	response, err := json.Marshal(highTemps)
	if err != nil {
		log.Fatal("err:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetLowTempsPerMonthHandler returns temperature lows for each month that has been querie
func (wc *WController) GetLowTempsPerMonthHandler(w http.ResponseWriter, r *http.Request) {
	lowTemps := wc.GetLowTemps(wc.DB.OpenConnection())

	response, err := json.Marshal(lowTemps)
	if err != nil {
		log.Fatal("err:", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetAverageTemperaturesByMonthHandler calls service to return average temperature by provided month
func (wc *WController) GetAverageTemperaturesByMonthHandler(w http.ResponseWriter, r *http.Request) {
	monthMap := make(map[string]int, 12)
	monthMap["january"] = 1
	monthMap["february"] = 2
	monthMap["march"] = 3
	monthMap["april"] = 4
	monthMap["may"] = 5
	monthMap["june"] = 6
	monthMap["july"] = 7
	monthMap["august"] = 8
	monthMap["sepetember"] = 9
	monthMap["october"] = 10
	monthMap["november"] = 11
	monthMap["december"] = 12

	month := chi.URLParam(r, "month")
	querymonth := monthMap[strings.ToLower(month)]

	averageTempByMonth, err := wc.GetAverageTempByMonth(querymonth, wc.DB.OpenConnection())
	if err != nil {
		log.Fatal("err", err.Error())
	}

	response, err := json.Marshal(averageTempByMonth)
	if err != nil {
		log.Fatal("err", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetDaysByWeatherHandler calls service to return all days grouped by weather description
func (wc *WController) GetDaysByWeatherHandler(w http.ResponseWriter, r *http.Request) {
	daysByWeather := wc.GetDaysByWeatherType(wc.DB.OpenConnection())

	response, err := json.Marshal(daysByWeather)
	if err != nil {
		log.Fatal("err", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
