package main

import (
	"fmt"

	"github.com/JumpNShootMan/RetoBCP/meme-bank/database"
	"github.com/JumpNShootMan/RetoBCP/meme-bank/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func setupRoutes(app *fiber.App) {

	app.Get("/api/v1/transactions", models.GetTransactions)
	app.Get("/api/v1/transaction/:id", models.GetTransaction)
	app.Post("/api/v1/transaction", models.NewTransaction)
	app.Delete("/api/v1/transaction/:id", models.DeleteTransaction)
	app.Get("/api/v1/transactions/user/:id", models.GetTransactionsOfUser)

	//////////////////////////////////////////////////////////
	app.Get("/api/v1/users", models.GetUsers)
	app.Get("/api/v1/user/:id", models.GetUser)
	app.Post("/api/v1/user", models.NewUser)
	app.Delete("/api/v1/user/:id", models.DeleteUser)
	//////////////////////////////////////////////////////////
	app.Post("/api/v1/login", models.LogIn)

	///////////////////////////////////////////////////////////
	app.Post("/api/v1/category", models.NewCategory)
	app.Get("/api/v1/categories", models.GetCategories)
	///////////////////////////////////////////////////////////
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/transactions", models.HandleWebsocket())
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "banks.db")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	database.DBConn.AutoMigrate(&models.Transaction{})
	database.DBConn.AutoMigrate(&models.User{})
	database.DBConn.AutoMigrate(&models.Category{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	AllowHeaders: "*",
	// }))

	setupRoutes(app)

	initDatabase()
	//models.CreateBank()
	app.Listen(":1996")

	defer database.DBConn.Close()
}
