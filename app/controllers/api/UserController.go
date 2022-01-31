package api

import (
	"github.com/anik4good/go_email/app/models"
	"github.com/anik4good/go_email/database"
	"github.com/gofiber/fiber/v2"
)

// var database *sql.DB

// Return all users as JSON
func GetAllUsers(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var Queues []models.QueuedEmail
		if response := db.Find(&Queues); response.Error != nil {
			panic("Error occurred while retrieving users from the database: " + response.Error.Error())
		}

		err := ctx.JSON(Queues)
		if err != nil {
			panic("Error occurred when returning JSON of users: " + err.Error())
		}
		return err
	}
}

// func CreateUser(c *fiber.Ctx) error {
// 	//	c.Send("Hello, World!")

// 	requestBody := c.Body()
// 	var email models.QueuedEmail
// 	json.Unmarshal(requestBody, &email)
// 	_, err := database.Exec(`INSERT INTO users(name, email,status) VALUES (?,?,?)`, email.Name, email.Email, email.Status)
// 	if err != nil {

// 		//	panic(err)

// 		fmt.Println("error creating user:", email.Name)
// 		//	json.NewEncoder(c).Encode("error creating user:")
// 		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 			"success": false,
// 			"data": fiber.Map{
// 				"todo": email,
// 			},
// 		})
// 		//	json.NewEncoder(c).Encode("error creating user:")
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"success": true,
// 		"data": fiber.Map{
// 			"todo": email,
// 		},
// 	})
// }
