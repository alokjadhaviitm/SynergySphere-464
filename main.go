package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
	"synergysphere464/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var store = session.New()

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
	sqlBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		log.Fatal("Failed to read init.sql file:", err)
	}
	sqlStatements := string(sqlBytes)

	// Execute the SQL statements
	_, err = db.Exec(sqlStatements)
	if err != nil {
		log.Fatal("Failed to execute init.sql statements:", err)
	}
	// Use html/template with Fiber
	engine := html.New("./templates", ".html")
	engine.Reload(true) // for dev: auto-reload templates

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static files
	app.Static("/static", "./static")

	// Route: Home
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layout", fiber.Map{
			"Title": "SynergySphere464",
		})
	})
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("signup", fiber.Map{
			"Title": "Sign Up - SynergySphere464",
		})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Login - SynergySphere464",
		})
	})

	app.Use("/dashboard", middleware.AuthSessionMiddleware(db))

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")

		if userID == nil {
			return c.Redirect("/login")
		}

		// Query user details from DB
		var firstName, email string
		query := `SELECT first_name, email FROM users WHERE id = $1`
		err := db.QueryRow(query, userID).Scan(&firstName, &email)
		if err != nil {
			log.Println("User fetch error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to load dashboard")
		}

		return c.Render("dashboard", fiber.Map{
			"Title":     "Dashboard - SynergySphere464",
			"UserName":  firstName,
			"UserEmail": email,
		})
	})

	// app.Get("/dashboard", func(c *fiber.Ctx) error {
	// 	return c.Render("dashboard", fiber.Map{
	// 		"Title": "dashboard - SynergySphere464",
	// 	})
	// })

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Get the session
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
		}

		// Destroy the session
		if err := sess.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Logout failed")
		}

		// Redirect to login
		return c.Redirect("/login")
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		firstName := strings.TrimSpace(c.FormValue("first_name"))
		lastName := strings.TrimSpace(c.FormValue("last_name"))
		email := strings.TrimSpace(c.FormValue("email"))
		password := c.FormValue("password")

		// Basic validation
		if firstName == "" || lastName == "" || email == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("All fields are required.")
		}
		if _, err := mail.ParseAddress(email); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid email.")
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Insert user into database
		query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)`
		_, err = db.Exec(query, firstName, lastName, email, string(hashedPassword))
		if err != nil {
			log.Println("DB insert error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Could not create account.")
		}

		// Redirect or respond
		return c.Redirect("/login")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		var id int
		var hashedPassword string
		var firstName string

		// Query user from DB
		err := db.QueryRow("SELECT id, password, first_name FROM users WHERE email = $1", email).Scan(&id, &hashedPassword, &firstName)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password.")
		}

		// Compare hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password.")
		}

		// Create session
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		sess.Set("user_id", id)
		sess.Set("user_name", firstName)
		sess.Save()

		return c.Redirect("/")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		userID := sess.Get("user_id")
		userName := sess.Get("user_name")

		if userID == nil {
			return c.Redirect("/login")
		}

		return c.Render("dashboard", fiber.Map{
			"Title": "Dashboard",
			"User":  userName,
		}) // your dashboard template
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
