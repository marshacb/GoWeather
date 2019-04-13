package weather

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Coordinates type for weather info
type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// Info primary weather data
type Info struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Weather list of weather data for a given day
type Weather []Info

// Main weather details
type Main struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
	MinTemp  float64 `json:"temp_min"`
	MaxTemp  float64 `json:"temp_max"`
}

// Wind information
type Wind struct {
	Speed  float64 `json:"speed"`
	Degree float64 `json:"deg"`
}

// Clouds weather
type Clouds struct {
	All int `json:"all"`
}

// OpenWeatherAPIData data captured from open weather api requests
type OpenWeatherAPIData struct {
	Location    string      `json:"location"`
	Coord       Coordinates `json:"coord"`
	Weather     Weather     `json:"weather"`
	Base        string      `json:"base"`
	MainWeather Main        `json:"main"`
	Visibility  float64     `json:"visibility"`
	Wind        Wind        `json:"wind"`
	Clouds      Clouds      `json:"clouds"`
}

// TemperatureRow date and temperature data types
type TemperatureRow struct {
	Date        time.Month `json:"month"`
	Temperature float64    `json:"temperature"`
}

// MainWeatherRow date and description
type MainWeatherRow struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// GetWeather calls openweather api with provided cityName, stores weather data city query value along with result in DB
func GetWeather(cityName string, weatherData OpenWeatherAPIData, db *sql.DB) OpenWeatherAPIData {
	_, err := storeQueryData(cityName, weatherData, db)
	if err != nil {
		fmt.Println("Error storing query:", err.Error())
	}
	weatherData.Location = cityName
	return weatherData
}

func storeQueryData(cityName string, weatherData OpenWeatherAPIData, db *sql.DB) (int, error) {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	query := `
		INSERT INTO WeatherQueries 
		(City, MainWeather, Description, Temperature, Pressure, MinTemp, MaxTemp, Latitude, Longitude, WeatherDate) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err = db.Query(
		query,
		cityName,
		weatherData.Weather[0].Main,
		weatherData.Weather[0].Description,
		weatherData.MainWeather.Temp,
		weatherData.MainWeather.Pressure,
		weatherData.MainWeather.MinTemp,
		weatherData.MainWeather.MaxTemp,
		weatherData.Coord.Latitude,
		weatherData.Coord.Longitude,
		time.Now(),
	)
	if err != nil {
		log.Fatal("err", err)
		return -1, err
	}
	return 1, nil
}

// GetQueryCount returns the total number of queries/records in DB
func GetQueryCount(db *sql.DB) (int, error) {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	query := `
		SELECT Count(*) from WeatherQueries
	`

	row, err := db.Query(query)
	if err != nil {
		fmt.Println("err", err)
	}
	var queryCount int64

	for row.Next() {
		err = row.Scan(&queryCount)
		if err != nil {
			return -1, err
		}
	}
	return int(queryCount), nil
}

// GetHighTemps returns high temps for all months found in queries table
func GetHighTemps(db *sql.DB) []TemperatureRow {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return nil
	}

	query := `
		 SELECT Month(WeatherDate), MAX(Temperature) 
		 FROM WeatherQueries WHERE YEAR(WeatherDate) = YEAR(GETDATE()) 
		 GROUP BY Month(WeatherDate)
	`

	maxTempsPerMonth := []TemperatureRow{}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("failed to query:", err.Error())
	}

	for rows.Next() {
		var tempRow TemperatureRow
		err = rows.Scan(&tempRow.Date, &tempRow.Temperature)
		if err != nil {
			log.Fatal("error scanning:", err.Error())
		}

		maxTempsPerMonth = append(maxTempsPerMonth, tempRow)
	}

	return maxTempsPerMonth
}

// GetLowTemps returns high temps for all months found in queries table
func GetLowTemps(db *sql.DB) []TemperatureRow {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return nil
	}

	query := `
		 SELECT Month(WeatherDate), MIN(Temperature) 
		 FROM WeatherQueries WHERE YEAR(WeatherDate) = YEAR(GETDATE()) 
		 GROUP BY Month(WeatherDate)
	`

	lowTempsPerMonth := []TemperatureRow{}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("failed to query:", err.Error())
	}

	for rows.Next() {
		var tempRow TemperatureRow
		err = rows.Scan(&tempRow.Date, &tempRow.Temperature)
		if err != nil {
			log.Fatal("error scanning:", err.Error())
		}

		lowTempsPerMonth = append(lowTempsPerMonth, tempRow)
	}

	fmt.Println("low temps per month", lowTempsPerMonth)

	return lowTempsPerMonth
}

// GetAverageTempByMonth returns the average temperatrue for a given month
func GetAverageTempByMonth(month int, db *sql.DB) (TemperatureRow, error) {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return TemperatureRow{}, err
	}

	query := `
		SELECT Month(WeatherDate), AVG(Temperature) 
		FROM WeatherQueries 
		Where MONTH(WeatherDate) = ` + strconv.Itoa(month) +
		`GROUP BY Month(WeatherDate)`

	row, err := db.Query(query)
	if err != nil {
		log.Fatal("err", err.Error())
	}

	var avgTempByMonth TemperatureRow
	for row.Next() {
		err = row.Scan(&avgTempByMonth.Date, &avgTempByMonth.Temperature)
		if err != nil {
			log.Fatal("err", err.Error())
		}
	}

	return avgTempByMonth, nil
}

// GetDaysByWeatherType returns all days in DB and their associated weather descriptions
func GetDaysByWeatherType(db *sql.DB) []MainWeatherRow {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return nil
	}

	query := `
		SELECT WeatherDate, CAST(MainWeather AS varchar(max)) 
		FROM WeatherQueries GROUP BY CAST(MainWeather AS varchar(max)), WeatherDate
	`

	daysByWeatherType := []MainWeatherRow{}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("failed to query:", err.Error())
	}

	for rows.Next() {
		var weatherRow MainWeatherRow
		err = rows.Scan(&weatherRow.Date, &weatherRow.Description)
		if err != nil {
			log.Fatal("error scanning:", err.Error())
		}

		daysByWeatherType = append(daysByWeatherType, weatherRow)
	}

	return daysByWeatherType
}
