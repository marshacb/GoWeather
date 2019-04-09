package bookmark

import (
	"context"
	"database/sql"
	"fmt"
)

// Location stores city as bookmarked location in DB
func Location(cityName string, db *sql.DB) string {
	_, err := storeLocationBookmark(cityName, db)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return "failed to bookmark"
	}

	return cityName + " bookmarked"
}

func storeLocationBookmark(cityName string, db *sql.DB) (int, error) {
	ctx := context.Background()
	defer db.Close()

	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	query := `
		USE WeatherTestDB;
		INSERT INTO BookmarkedLocations 
		(CityName, Bookmarked) 
		VALUES (?, ?);
	`
	_, err = db.Query(query, cityName, 1)
	if err != nil {
		fmt.Println("err", err.Error())
		return -1, err
	}

	return 1, nil
}
