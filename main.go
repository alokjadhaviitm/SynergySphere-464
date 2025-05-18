package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Setup DB
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "synergydb"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connect error:", err)
	}
	defer db.Close()

	// Use html/template with Fiber
	engine := html.New("templates", ".html")
	engine.Reload(true) // for dev: auto-reload templates

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static files
	app.Static("/static", "./static")

	// Route: Home
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layout", fiber.Map{
			"Title": "SynergySphere",
		})
	})

	log.Println("Fiber server running at http://localhost:8070")
	log.Fatal(app.Listen(":8070"))
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
