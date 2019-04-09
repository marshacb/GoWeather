package models

import "time"

// WeatherData : holds data retrieved from openweatherapi requests
type WeatherData struct {
	QueryString string    `json:"query_string"`
	CityName    string    `json:"city_name"`
	WeatherType string    `json:"weather_type"`
	Description string    `json:"description"`
	Temperature float64   `json:"temperature"`
	Pressure    float64   `json:"pressure"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}
