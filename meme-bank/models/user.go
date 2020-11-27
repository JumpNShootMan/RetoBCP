package models

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JumpNShootMan/RetoBCP/meme-bank/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string  `json:"name" validate:"required"`
	Surname  string  `json:"surname" validate:"required"`
	Email    string  `json:"email" validate:"required"`
	Password string  `json:"-" validate:"required"`
	Balance  float64 `json:"balance" gorm:"-"`
}

// LoginRequest is the format to login a user
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Balance returns the balance for an account by adding all transactions
func Balance(userID uint) (float64, error) {
	db := database.DBConn
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return 0, fmt.Errorf("user not found")
	}
	transactions, err := TransactionsOfUser(userID)
	if err != nil {
		return 0, fmt.Errorf("no transactions found for user")
	}

	balance := 0.0

	for _, transaction := range transactions {
		if transaction.FromID == userID {
			balance = balance - transaction.Amount
		}
		if transaction.ToID == userID {
			balance = balance + transaction.Amount
		}
	}

	return balance, nil
}

// GetUsers ...
func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var user []User
	db.Find(&user)
	return c.JSON(user)
}

// GetUser ...
func GetUser(c *fiber.Ctx) error {
	idparam := c.Params("id")
	db := database.DBConn
	var user User
	if err := db.First(&user, idparam).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No user found with given id"})

	}
	id, _ := strconv.ParseUint(idparam, 10, 32)
	balance, err := Balance(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error calculating balance"})
	}
	user.Balance = balance
	return c.Status(fiber.StatusOK).JSON(user)
}

// NewUser ...
func NewUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if err := Validator().Struct(user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"validation error": err.Error()})
	}

	db.Create(&user)
	return c.Status(fiber.StatusOK).JSON(user)
}

// DeleteUser ...
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var user User
	db.First(&user, id)
	if user.Name == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No user found with given id"})
	}
	db.Delete(&user)
	return c.Status(http.StatusNoContent).JSON("user deleted")
}

// LogIn ...
func LogIn(c *fiber.Ctx) error {
	// TODO: Change to more secure login
	db := database.DBConn
	login := new(LoginRequest)
	if err := c.BodyParser(login); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(err)
	}

	if err := Validator().Struct(login); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"validation error": err.Error()})
	}

	var user User
	if err := db.Where("Email = ? AND Password = ?", login.Email, login.Password).First(&user).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect email or password"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"user": user.ID})
}
