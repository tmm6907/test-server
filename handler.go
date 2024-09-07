package main

import (
	"database/sql"

	"github.com/labstack/echo"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) SetTitle(c echo.Context, val string) {
	c.Response().Header().Set("X-Page-Title", val)
}
