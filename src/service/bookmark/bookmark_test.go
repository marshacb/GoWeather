package bookmark_test

import (
	"errors"
	"go_weather/src/service/bookmark"
	"log"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bookmark", func() {
	Context("Location", func() {
		It("successfully calls DB to store/bookmark provided location", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			cityToBookmark := "Charlotte"
			expectedQuery := "USE WeatherTestDB; INSERT INTO BookmarkedLocations \\(CityName, Bookmarked\\) VALUES \\(\\?, \\?\\);"
			mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"1", "1"}))
			mock.ExpectBegin()

			response := bookmark.Location(cityToBookmark, db)
			Expect(response).To(Equal(cityToBookmark + " bookmarked"))
		})

		It("returns an error when call to DB fails", func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			cityToBookMark := "Raleigh"
			connectError := errors.New("failed to connect")
			expectedQuery := "USE WeatherTestDB; INSERT INTO BookmarkedLocations \\(CityName, Bookmarked\\) VALUES \\(\\?, \\?\\);"
			mock.ExpectQuery(expectedQuery).WithArgs(cityToBookMark, 1).WillReturnError(connectError)
			mock.ExpectBegin()

			response := bookmark.Location(cityToBookMark, db)
			Expect(response).To(Equal("failed to bookmark"))
		})
	})
})
