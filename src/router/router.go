package router

import (
	"net/http"
	"time"

	"../controllers/bookmarks"
	"../controllers/weathers"

	"../db"
	"../service/bookmark"
	"../service/weather"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Initialize chi mux router
func Initialize() *chi.Mux {
	mssqlDB := &db.MSSQLDB{}
	bookmarksController := bookmarks.BookMarkController{DB: mssqlDB, Location: bookmark.Location}
	weatherController := weathers.WController{
		DB:                    mssqlDB,
		GetWeather:            weather.GetWeather,
		GetQueryCount:         weather.GetQueryCount,
		GetHighTemps:          weather.GetHighTemps,
		GetLowTemps:           weather.GetLowTemps,
		GetAverageTempByMonth: weather.GetAverageTempByMonth,
		GetDaysByWeatherType:  weather.GetDaysByWeatherType,
		GetOpenWeatherAPIData: http.Get,
	}
	muxRouter := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
	muxRouter.Use(cors.Handler)
	muxRouter.Use(middleware.RequestID)
	muxRouter.Use(middleware.RealIP)
	muxRouter.Use(middleware.Logger)
	muxRouter.Use(middleware.Recoverer)
	muxRouter.Use(middleware.Timeout(200 * time.Second))

	muxRouter.Get("/get-city-weather/{city}", weatherController.GetWeatherHandler)
	muxRouter.Get("/total-queries", weatherController.GetQueryCountHandler)
	muxRouter.Get("/high-temperatures", weatherController.GetHighTempsPerMonthHandler)
	muxRouter.Get("/low-temperatures", weatherController.GetLowTempsPerMonthHandler)
	muxRouter.Get("/average-temperature/{month}", weatherController.GetAverageTemperaturesByMonthHandler)
	muxRouter.Get("/weather-days", weatherController.GetDaysByWeatherHandler)
	muxRouter.Post("/bookmark-location", bookmarksController.BookmarkLocationHandler)
	return muxRouter
}
