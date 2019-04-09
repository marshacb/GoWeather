package weathers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go_weather/controllers/weathers"
	"go_weather/service/weather"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mssqlTestDB struct {
}

func (m *mssqlTestDB) OpenConnection() *sql.DB {
	db, _, _ := sqlmock.New()
	return db
}

var _ = Describe("Weathers", func() {
	Context("GetWeatherHandler", func() {
		It("returns the correct weather data for the specified location", func() {
			location := "Harlem"
			weatherData := weather.OpenWeatherAPIData{
				Location:    location,
				Coord:       weather.Coordinates{Latitude: 11.7, Longitude: 10.5},
				Weather:     weather.Weather{weather.Info{ID: 10, Main: "Cloudy", Description: "mostly cloudy", Icon: "04b"}},
				Base:        "stations",
				MainWeather: weather.Main{Temp: 99.5, Pressure: 100, Humidity: 100, MinTemp: 87.6, MaxTemp: 100.2},
				Visibility:  15.7,
				Wind:        weather.Wind{Speed: 20.0, Degree: 10.0},
				Clouds:      weather.Clouds{All: 40},
			}
			body, _ := json.Marshal(weatherData)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/get-city-weather/"+location, bytes.NewReader(body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("city", location)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			msDB := &mssqlTestDB{}

			getWeather := func(string, weather.OpenWeatherAPIData, *sql.DB) weather.OpenWeatherAPIData {
				return weatherData
			}
			getQueryCount := func(*sql.DB) (int, error) {
				return 1, nil
			}
			getHighTemps := func(*sql.DB) []weather.TemperatureRow {
				return []weather.TemperatureRow{}
			}
			getLowTemps := func(*sql.DB) []weather.TemperatureRow {
				return []weather.TemperatureRow{}
			}

			getAverageTempByMonth := func(int, *sql.DB) (weather.TemperatureRow, error) {
				return weather.TemperatureRow{}, nil
			}

			getDaysByWeatherType := func(*sql.DB) []weather.MainWeatherRow {
				return []weather.MainWeatherRow{}
			}

			getOpenWeatherApiData := func(url string) (resp *http.Response, err error) {
				recorder := httptest.ResponseRecorder{}
				return recorder.Result(), nil
			}

			weatherController := weathers.WController{
				DB:                    msDB,
				GetWeather:            getWeather,
				GetQueryCount:         getQueryCount,
				GetHighTemps:          getHighTemps,
				GetLowTemps:           getLowTemps,
				GetAverageTempByMonth: getAverageTempByMonth,
				GetDaysByWeatherType:  getDaysByWeatherType,
				GetOpenWeatherAPIData: getOpenWeatherApiData,
			}

			handler := http.HandlerFunc(weatherController.GetWeatherHandler)
			handler.ServeHTTP(w, r)

			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Location"}).
				AddRow("Harlem").
				AddRow("Bronx").
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("INSERT").WillReturnRows(rows)
			mock.ExpectBegin()

			var response weather.OpenWeatherAPIData

			err = json.NewDecoder(r.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.Location).To(Equal(location))
		})
	})
})
