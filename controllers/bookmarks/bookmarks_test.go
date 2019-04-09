package bookmarks_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go_weather/controllers/bookmarks"
	"log"
	"net/http"
	"net/http/httptest"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mssqlTestDB struct {
}

func (m *mssqlTestDB) OpenConnection() *sql.DB {
	db, _, _ := sqlmock.New()
	return db
}

var _ = Describe("Bookmarks", func() {
	Context("BookmarkLocationHandler", func() {
		It("returns string based on requested bookmarked city", func() {
			postBookmarkDataBody := struct {
				CityName string
			}{
				"Oakland",
			}
			body, _ := json.Marshal(postBookmarkDataBody)

			req, _ := http.NewRequest("POST", "/bookmark-location", bytes.NewReader(body))
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("err", err)
			}
			defer db.Close()

			expectedQuery := "USE WeatherTestDB; INSERT INTO BookmarkedLocations \\(CityName, Bookmarked\\) VALUES \\(\\?, \\?\\);"
			mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"1", "1"}))
			mock.ExpectBegin()

			testLocation := func(string, *sql.DB) string {
				return "Oakland bookmarked"
			}

			msDB := &mssqlTestDB{}

			bookMarkController := bookmarks.BookMarkController{DB: msDB, Location: testLocation}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(bookMarkController.BookmarkLocationHandler)
			handler.ServeHTTP(recorder, req)

			type PostBookMarkResponse struct {
				Message string `json:"message"`
			}

			var response PostBookMarkResponse

			err = json.NewDecoder(recorder.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.Message).To(Equal("Oakland bookmarked"))
		})
	})
})
