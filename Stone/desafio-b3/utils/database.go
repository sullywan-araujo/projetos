package utils

import (
	"log"

	"desafio-b3/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func ConectarBanco() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:cotacoes.db?cache=shared&mode=rwc"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	db.AutoMigrate(&models.Negocio{})
	return db
}
