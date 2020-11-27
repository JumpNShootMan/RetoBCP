package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JumpNShootMan/RetoBCP/meme-bank/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

// Transaction - modelo
type Transaction struct {
	gorm.Model
	Amount      float64  `json:"amount" validate:"required"`
	Category    Category `json:"category" validate:"-"`
	CategoryID  uint     `json:"category_id" validate:"-"`
	FromID      uint     `json:"from_user_id" validate:"required"`
	From        User     `validate:"-"`
	ToID        uint     `json:"to_user_id" validate:"required"`
	To          User     `validate:"-"`
	Description string   `json:"description" validate:"required"`
}

// TransactionsOfUser retorna todas las transacciones del usuario
func TransactionsOfUser(userID uint) ([]Transaction, error) {
	db := database.DBConn
	if err := db.First(&User{}, userID).Error; err != nil {
		return nil, fmt.Errorf("user does not exist")
	}

	var transactions []Transaction
	if err := db.Preload("From").Preload("To").Preload("Category").Order("created_at desc").Where("from_id = ? OR to_id = ?", userID, userID).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("no transactions found for that user")
	}

	return transactions, nil
}

// GetTransactionsOfUser consigue todas las transacciones del usuario
func GetTransactionsOfUser(c *fiber.Ctx) error {
	idparam := c.Params("id")
	id, err := strconv.ParseUint(idparam, 10, 16)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user does not exist"})
	}

	transactions, err := TransactionsOfUser(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user or transaction does not exist"})
	}

	return c.Status(fiber.StatusOK).JSON(transactions)
}

// GetTransactions ...
func GetTransactions(c *fiber.Ctx) error {
	db := database.DBConn
	var transactions []Transaction
	db.Preload("From").Preload("To").Preload("Category").Find(&transactions)
	return c.Status(fiber.StatusOK).JSON(transactions)
}

// GetTransaction ...
func GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var transaction Transaction
	if err := db.Find(&transaction, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No transaction found with given id"})
	}
	return c.Status(fiber.StatusOK).JSON(transaction)
}

// NewTransaction ...
func NewTransaction(c *fiber.Ctx) error {

	db := database.DBConn
	transaction := new(Transaction)
	if err := c.BodyParser(transaction); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if err := Validator().Struct(transaction); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"validation error": err.Error()})
	}

	userTo := new(User)
	userFrom := new(User)
	category := new(Category)

	if err := db.First(&userTo, transaction.ToID).Error; err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "to_user does not exist"})
	}

	if err := db.First(&userFrom, transaction.FromID).Error; err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "from_user does not exist"})
	}

	if err := db.First(&category, transaction.CategoryID).Error; err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "category does not exist"})
	}

	userToBalance, err := Balance(transaction.FromID)
	if transaction.Amount > userToBalance {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "not enough funds to complete transaction"})
	}

	db.Create(&transaction)
	db.Preload("To").Preload("From").Preload("Category").First(&transaction, transaction.ID)

	marshal, err := json.Marshal(transaction)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}

	Publish(marshal)
	return c.Status(fiber.StatusOK).JSON(transaction)
}

// DeleteTransaction ...
func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var transaction Transaction

	if err := db.First(&transaction, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No transaction found with given id"})
	}

	db.Delete(&transaction)
	return c.Status(http.StatusNoContent).JSON(fiber.Map{"message": "transaction deleted"})
}
