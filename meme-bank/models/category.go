package models

import (
	"net/http"

	"github.com/JumpNShootMan/RetoBCP/meme-bank/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

// Category model
type Category struct {
	gorm.Model
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// GetCategories ...
func GetCategories(c *fiber.Ctx) error {
	db := database.DBConn
	var cat []Category
	db.Find(&cat)
	return c.JSON(cat)
}

// GetCategory ...
func GetCategory(c *fiber.Ctx) error {
	idparam := c.Params("id")
	db := database.DBConn
	var cat Category
	if err := db.First(&cat, idparam).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No category found with given id"})

	}

	return c.Status(fiber.StatusOK).JSON(cat)
}

// NewCategory ...
func NewCategory(c *fiber.Ctx) error {
	db := database.DBConn
	cat := new(Category)
	if err := c.BodyParser(cat); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if err := Validator().Struct(cat); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"validation error": err.Error()})
	}

	db.Create(&cat)
	return c.Status(fiber.StatusOK).JSON(cat)
}

// DeleteCategory ...
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var cat Category
	db.First(&cat, id)
	if cat.Name == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No category found with given id"})
	}
	db.Delete(&cat)
	return c.Status(http.StatusNoContent).JSON("category deleted")
}
