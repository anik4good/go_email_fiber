package routes

import (
	Controller "github.com/anik4good/go_email/app/controllers/api"
	"github.com/anik4good/go_email/database"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(api fiber.Router, db *database.Database) {
	// registerRoles(api, db)
	registerUsers(api, db)
}

func registerUsers(api fiber.Router, db *database.Database) {
	users := api.Group("/users")

	users.Get("/", Controller.GetAllUsers(db))
	// users.Get("/:id", Controller.GetUser(db))
	// users.Post("/", Controller.AddUser(db))
	// users.Put("/:id", Controller.EditUser(db))
	// users.Delete("/:id", Controller.DeleteUser(db))
}
