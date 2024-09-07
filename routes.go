package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) Index(c echo.Context) error {
	h.SetTitle(c, "Home Page")
	data := map[string]any{
		"Selected": 1,
		"Title":    "Home Page",
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) CalendarPage(c echo.Context) error {
	h.SetTitle(c, "Test Page")
	data := map[string]any{
		"Selected": 0,
		"Title":    "Calendar Page",
	}
	return c.Render(http.StatusOK, "calendar.html", data)
}

func (h *Handler) GetEvents(c echo.Context) error {
	var results []map[string]any // Assuming you want to collect results in a slice
	rows, err := h.DB.Query("SELECT * FROM events")
	if err != nil {
		c.Error(err)
		log.Println(err)
		return nil
	}
	defer rows.Close() // Ensure rows are closed after processing

	for rows.Next() {
		var result map[string]any
		err = rows.Scan(&result) // Adjust based on your actual schema
		if err != nil {
			c.Error(err)
			log.Println(err)
			return nil
		}
		results = append(results, result)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.Error(err)
		log.Println(err)
		return nil
	}

	err = c.JSON(http.StatusOK, results)
	if err != nil {
		return err
	}
	return nil
}
