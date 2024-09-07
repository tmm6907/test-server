package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := fmt.Errorf("Template not found '%s'", name)
		slog.Error(err.Error())
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

func initDB() (*sql.DB, error) {
	// Open a connection to the SQLite database file (or create it if it doesn't exist)
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create a sample table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		startDate TEXT,
		startTime TEXT,
		endDate TEXT,
		endTime TEXT,
		description TEXT,
		created_at NOT NULL DEFAULT current_timestamp,
		updated_at NOT NULL DEFAULT current_timestamp
	);
	`
	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

func populateTemplates(pattern string) (map[string]*template.Template, error) {
	templs := make(map[string]*template.Template)
	files, err := filepath.Glob(pattern)
	layoutFiles := []string{"templates/base.html", "templates/styles.html"}
	if err != nil {
		return templs, nil
	}
	mainTemplate, err := template.ParseFiles(layoutFiles...)
	if err != nil {
		log.Fatalf("Error parsing main layout: %v", err)
	}
	for _, file := range files {
		if slices.Contains(layoutFiles, file) {
			continue
		}
		fileName := filepath.Base(file)

		clonedTemplate, err := mainTemplate.Clone()
		if err != nil {
			log.Fatalf("Error cloning main template: %v", err)
		}
		templs[fileName] = template.Must(clonedTemplate.ParseFiles(file))
		slog.Info(fmt.Sprintf("Successfully parsed template: %s", fileName))
	}
	return templs, nil
}

func main() {
	server := echo.New()
	templates, err := populateTemplates("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
	t := &Template{
		templates: templates,
	}
	server.Renderer = t
	// server.Logger.SetOutput(os.Stdout)

	// // Add middleware to log errors
	// server.Use(middleware.Logger())
	// server.Use(middleware.Recover())

	h := &Handler{
		DB: nil,
	}
	server.GET("/", h.Index)
	server.GET("/events", h.GetEvents)
	server.GET("/test", h.CalendarPage)
	log.Println("Registered routes:")

	port := "8000"
	log.Printf("listening on port %s\n", port)
	log.Fatal(server.Start(fmt.Sprintf(":%s", port)))
}
