package models

import (
	"github.com/jamilchioino/meme-bank/database"
)

// CreateBank creates an initial transaction from a reserve bank
func CreateBank() error {
	db := database.DBConn
	bank := User{
		Email:    "reservas@interbank.com",
		Name:     "Interbank",
		Password: "12345",
		Surname:  "Peru",
	}

	reserves := User{
		Email:    "banco@nacional.com",
		Name:     "Reservas",
		Surname:  "Nacionales",
		Password: "12345",
	}

	if err := db.Create(&bank).Error; err != nil {
		return err
	}
	if err := db.Create(&reserves).Error; err != nil {
		return err
	}

	transaction := Transaction{
		Amount:      10000,
		Description: "Emisión de Billetes",
		FromID:      reserves.ID,
		ToID:        bank.ID,
	}

	if err := db.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}
