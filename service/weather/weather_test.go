package weather_test

import (
	"fmt"
	"go_weather/service/weather"
	"log"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Weather", func() {
	Context("GetWeather", func() {
		It("Successfully calls DB to store provided weather query data", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			weatherData := weather.OpenWeatherAPIData{
				"dallas",
				weather.Coordinates{Latitude: 11.7, Longitude: 10.5},
				weather.Weather{weather.Info{ID: 10, Main: "Cloudy", Description: "mostly cloudy", Icon: "04b"}},
				"stations",
				weather.Main{Temp: 99.5, Pressure: 100, Humidity: 100, MinTemp: 87.6, MaxTemp: 100.2},
				15.7,
				weather.Wind{Speed: 20.0, Degree: 10.0},
				weather.Clouds{All: 40},
			}

			cityName := "dallas"
			expectedQuery := "INSERT INTO WeatherQueries \\(City, MainWeather, Description, Temperature, Pressure, MinTemp, MaxTemp, Latitude, Longitude, WeatherDate\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\);"
			mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"houston", "Cloudy"}))
			mock.ExpectBegin()

			response := weather.GetWeather(cityName, weatherData, db)
			Expect(response.Location).To(Equal(cityName))
			Expect(response.Weather[0].Main).To(Equal(weatherData.Weather[0].Main))
		})
	})

	Context("GetQueryCount", func() {
		It("Successfully calls DB to return total number of queries entered", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Count"}).
				AddRow(10).
				AddRow(1).
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectBegin()

			response, _ := weather.GetQueryCount(db)
			Expect(int(response)).To(Equal(10))
		})
	})

	Context("GetHighTemps", func() {
		It("Successfully calls DB to return high temps by month", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Date", "Temperature"}).
				AddRow(5, 100).
				AddRow(4, 50).
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectBegin()

			response := weather.GetHighTemps(db)

			Expect(response[0].Date).To(Equal(time.Month(5)))
			Expect(response[0].Temperature).To(Equal(float64(100)))
		})
	})

	Context("GetLowTemps", func() {
		It("Successfully calls DB to return low temps by month", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Date", "Temperature"}).
				AddRow(10, 50).
				AddRow(4, 2).
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectBegin()

			response := weather.GetLowTemps(db)

			Expect(response[0].Date).To(Equal(time.Month(10)))
			Expect(response[0].Temperature).To(Equal(float64(50)))
		})
	})

	Context("GetAverageTempByMonth", func() {
		It("Successfully calls DB to return avg temp for the given month", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Month", "Temperature"}).
				AddRow(2, 25).
				AddRow(4, 2).
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectBegin()

			month := 2
			response, err := weather.GetAverageTempByMonth(month, db)
			if err != nil {
				fmt.Println("err", err)
			}

			Expect(response.Date).To(Equal(time.Month(2)))
			Expect(response.Temperature).To(Equal(float64(25)))
		})
	})

	Context("GetDaysByWeatherType", func() {
		It("Successfully calls DB to return all days by weather type", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"Date", "Description"}).
				AddRow(time.Now(), "Cloudy").
				AddRow(4, 2).
				RowError(1, fmt.Errorf("row error"))
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectBegin()

			response := weather.GetDaysByWeatherType(db)
			if err != nil {
				fmt.Println("err", err)
			}

			Expect(response[0].Description).To(Equal("Cloudy"))
		})
	})

})
