package main

import (
	"log"
	"net/http"
)

func (h Handler) GetEvents() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var results map[string]any
		rows, err := h.DB.Query("SELECT * FROM events")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(results)
			if err != nil {
				panic(err)
			}
			log.Println(results)
		}
	}
}
