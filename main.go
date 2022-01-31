package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/anik4good/go_email/app/models"
	"github.com/anik4good/go_email/config"
	configuration "github.com/anik4good/go_email/config"
	"github.com/anik4good/go_email/database"
	"github.com/anik4good/go_email/routes"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/session"
	// hashing "github.com/thomasvvugt/fiber-hashing"
)

type App struct {
	*fiber.App

	DB *database.Database
}

var sql_db *sql.DB
var logger *log.Logger

func main() {

	logger = config.InitLogger()
	config := configuration.New()

	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}
	// Initialize database
	db, err := database.New(&database.DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER"),
		Host:     config.GetString("DB_HOST"),
		Username: config.GetString("DB_USERNAME"),
		Password: config.GetString("DB_PASSWORD"),
		Port:     config.GetInt("DB_PORT"),
		Database: config.GetString("DB_DATABASE"),
	})

	// Auto-migrate database models
	if err != nil {
		fmt.Println("failed to connect to database:", err.Error())
	} else {
		if db == nil {
			fmt.Println("failed to connect to database: db variable is nil")
		} else {
			app.DB = db
			err := app.DB.AutoMigrate(&models.QueuedEmail{})
			if err != nil {
				fmt.Println("failed to automigrate role model:", err.Error())
				return
			}
			// err = app.DB.AutoMigrate(&models.User{})
			// if err != nil {
			// 	fmt.Println("failed to automigrate user model:", err.Error())
			// 	return
			// }
		}
	}

	// Register application API routes (using the /api/v1 group)
	api := app.Group("/api")
	apiv1 := api.Group("/v1")
	routes.RegisterAPI(apiv1, app.DB)

	// Custom 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		if err := c.SendStatus(fiber.StatusNotFound); err != nil {
			panic(err)
		}
		if err := c.Render("errors/404", fiber.Map{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return err
	})

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()

	// Start listening on the specified address
	err = app.Listen(config.GetString("APP_ADDR"))
	if err != nil {
		app.exit()
	}

	// do a forever for loop here that repetedly query the database and check for any record that has
	for {
		//	UserSeed()

		newRecords := checkForNewRecords()

		for newRecords.Next() {
			var q models.QueuedEmail
			err := newRecords.Scan(&q.ID, &q.Name, &q.Email)
			if err != nil {
				logger.Println("Error writting new records to queued sms struct", err)
				continue
			}

			changeStatusToPending(q.ID)

			go processEmail(q)
			logger.Println("Status change for id ", q.ID)

			fmt.Println("Status change for id ", q.ID)

		}

		time.Sleep(2 * time.Second)
		fmt.Println("No data found")
	}
}

func processEmail(queuedEmail models.QueuedEmail) {
	err := sendEmail(queuedEmail)
	if err != nil {
		logger.Println(err)
		return
	}

	changeStatusToSuccess(queuedEmail.ID)

}

func sendEmail(queuedEmail models.QueuedEmail) error {
	//	logger.Println("Sending sms to", queuedEmail.Email)
	fmt.Println("Sending Email to", queuedEmail.Email)

	//	send_email(queuedEmail)

	return nil
}

func checkForNewRecords() *sql.Rows {
	rows, err := sql_db.Query("select id, name, email from users WHERE status = 0 LIMIT 500")
	if err != nil {
		logger.Println("Error on new records checking ..", err)
	}
	return rows
}

func changeStatusToPending(id uint64) {
	_, err := sql_db.Exec("UPDATE users SET status = ? WHERE id = ?", 2, id)
	if err != nil {
		logger.Println("Error updating status of "+string(rune(id))+" in users", (err), " to pending")
		return
	}
}

func changeStatusToSuccess(id uint64) {
	_, err := sql_db.Exec("UPDATE users SET status = ? WHERE id = ?", 3, id)
	if err != nil {
		logger.Println("Error updating status of "+string(rune(id))+" in users", (err), " to pending")
		return
	}
}

// func send_email(queuedEmail models.QueuedEmail) {
// 	// Configuration
// 	from := "anik4nobody@gmail.com"
// 	password := "kfxmwyvzhcifoylj"
// 	//	to := []queuedEmail.Email
// 	to := []string{queuedEmail.Email}
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	//message := []byte("My super secret message.")

// 	message := []byte(
// 		"Subject: discount Gophers!\r\n" +
// 			"\r\n" +
// 			"This is the email body.\r\n")

// 	// Create authentication
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Send actual message
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Message sent to: ", queuedEmail.Email)

// }

func UserSeed() {

	for i := 0; i < 100; i++ {
		//prepare the statement
		//	stmt, _ := s.db.Prepare(`INSERT INTO users(name, email) VALUES (?,?)`)
		// execute query
		//	_, err := stmt.Exec(faker.Name(), faker.Email())

		_, err := sql_db.Exec(`INSERT INTO users(name, email,status) VALUES (?,?,?)`, faker.Name(), faker.Email(), 0)
		if err != nil {
			panic(err)
		}

	}
}

// Stop the Fiber application
func (app *App) exit() {
	_ = app.Shutdown()
}
