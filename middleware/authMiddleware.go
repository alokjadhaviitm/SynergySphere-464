package middleware

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New() // Or inject if already created

func AuthSessionMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get session
		sess, err := store.Get(c)
		if err != nil {
			log.Println("Session error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}

		// Check user ID from session
		userID := sess.Get("user_id")
		if userID == nil {
			return c.Redirect("/login")
		}

		// Check if user exists in DB
		var exists bool
		query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
		err = db.QueryRow(query, userID).Scan(&exists)
		if err != nil {
			log.Println("DB check error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("User validation failed")
		}

		if !exists {
			return c.Status(fiber.StatusUnauthorized).SendString("User not found")
		}

		// Set user ID in request context
		c.Locals("user_id", userID)

		return c.Next()
	}
}
