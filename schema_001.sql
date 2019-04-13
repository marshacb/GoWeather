CREATE DATABASE WeatherTestDB;

USE WeatherTestDB; 

CREATE TABLE BookmarkedLocations (
	ID INT NOT NULL IDENTITY PRIMARY KEY,
    CityName TEXT,
	BookMarked BIT
)

CREATE TABLE WeatherQueries (
	ID INT NOT NULL IDENTITY PRIMARY KEY,
	City TEXT,
	MainWeather TEXT,
	Description TEXT,
	Temperature DECIMAL,
	Pressure DECIMAL,
	MinTemp DECIMAL,
	MaxTemp DECIMAL,
	Longititude DECIMAL,
	Latitude DECIMAL,
	WeatherDate DATETIME
)