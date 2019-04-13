Getting Started
========================
Clone Project - "git clone https://github.com/marshacb/GoWeather.git"   
Run "docker-compose up" from root directory of project to create docker container  
Access api from localhost:8080   

Routes
==============
GET "/get-city-weather/{city}" - returns the weather data for the provided city/location   
GET "/total-queries" - returns the total number of queries   
GET "/high-temperatures" - returns the high temps for each month   
GET "/low-temperatures" - returns the low temps for each month   
GET "/average-temperature/{month}" - returns the average temp for the provided month   
GET "/weather-days" - returns all days queried grouped by weather description   
POST "/bookmark-location" - request body sample {"cityName": "Miami"}   


Testing
=============
To run Unit tests locally, install ginkgo and gomega on your machine by running 'go get https://github.com/onsi/ginkgo' and 'go get https://github.com/onsi/gomega' respectively.
Once installed run the gingko -r command inside of the go_weather/src directory.
Note: file paths may need to be updated due to dockerizing   