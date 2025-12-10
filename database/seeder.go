package database

import (
	"catatan-keuangan/modules/categories"
	"log"

	"gorm.io/gorm"
)

func seedCategories(db *gorm.DB) {
	categories := []categories.Category{
		{Name: "Makanan", Type: "expense"},
		{Name: "Transportasi", Type: "expense"},
		{Name: "Hiburan", Type: "expense"},
		{Name: "Kesehatan", Type: "expense"},
		{Name: "Pendidikan", Type: "expense"},
		{Name: "Tagihan", Type: "expense"},
		{Name: "Belanja", Type: "expense"},
		{Name: "Investasi", Type: "expense"},
		{Name: "Tabungan", Type: "expense"},
		{Name: "Lain-lain", Type: "expense"},
		{Name: "Gaji", Type: "income"},
		{Name: "Bonus", Type: "income"},
		{Name: "Hadiah", Type: "income"},
		{Name: "Penjualan", Type: "income"},
		{Name: "Pendapatan Sampingan", Type: "income"},
		{Name: "Lain-lain", Type: "income"},
	}
	for _, category := range categories {
		if err := db.Where("name = ? AND type = ?", category.Name, category.Type).
			FirstOrCreate(&category).Error; err != nil {
			log.Printf("gagal membuat kategori %s: %v\n", category.Name, err)
		}
	}
}
