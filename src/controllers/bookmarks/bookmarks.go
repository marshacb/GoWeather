package bookmarks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_weather/src/db"
	"net/http"
)

type locationToBookmark struct {
	CityName string `json:"cityName"`
}

// BookMarkController with dependencies
type BookMarkController struct {
	DB       db.DatabaseInterf
	Location func(string, *sql.DB) string
}

// BookmarkLocationHandler calls bookmark service to store location
func (b *BookMarkController) BookmarkLocationHandler(w http.ResponseWriter, r *http.Request) {
	var location locationToBookmark

	json.NewDecoder(r.Body).Decode(&location)
	bookmarked := b.Location(location.CityName, b.DB.OpenConnection())

	responseStruct := struct {
		Message string `json:"message"`
	}{
		bookmarked,
	}

	response, err := json.Marshal(responseStruct)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
