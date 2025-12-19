package database

import (
	"catatan-keuangan/modules/categories"
	"log"

	"gorm.io/gorm"
)

func stringPtr(s string) *string {
	return &s
}

func seedCategories(db *gorm.DB) {
	categories := []categories.Category{
		{
			Name:  "Gaji Bulanan",
			Type:  "income",
			Color: stringPtr("#10B981"),
			Icon:  "TrendingUp",
		},
		{
			Name:  "Freelance Project",
			Type:  "income",
			Color: stringPtr("#3B82F6"),
			Icon:  "Briefcase",
		},
		{
			Name:  "Bonus & Tunjangan",
			Type:  "income",
			Color: stringPtr("#8B5CF6"),
			Icon:  "Gift",
		},
		{
			Name:  "Investasi",
			Type:  "income",
			Color: stringPtr("#F59E0B"),
			Icon:  "TrendingUp",
		},
		{
			Name:  "Kopi & Snack",
			Type:  "expense",
			Color: stringPtr("#F59E0B"),
			Icon:  "Coffee",
		},
		{
			Name:  "Makan Siang",
			Type:  "expense",
			Color: stringPtr("#EF4444"),
			Icon:  "Utensils",
		},
		{
			Name:  "Belanja Bulanan",
			Type:  "expense",
			Color: stringPtr("#A855F7"),
			Icon:  "ShoppingBag",
		},
		{
			Name:  "Fashion & Pakaian",
			Type:  "expense",
			Color: stringPtr("#EC4899"),
			Icon:  "Shirt",
		},
		{
			Name:  "Listrik & Air",
			Type:  "expense",
			Color: stringPtr("#3B82F6"),
			Icon:  "Home",
		},
		{
			Name:  "Internet & Telepon",
			Type:  "expense",
			Color: stringPtr("#06B6D4"),
			Icon:  "Wifi",
		},
		{
			Name:  "Transport",
			Type:  "expense",
			Color: stringPtr("#64748B"),
			Icon:  "Car",
		},
		{
			Name:  "Bensin & Parkir",
			Type:  "expense",
			Color: stringPtr("#78716C"),
			Icon:  "Fuel",
		},
		{
			Name:  "Hiburan",
			Type:  "expense",
			Color: stringPtr("#F97316"),
			Icon:  "Popcorn",
		},
		{
			Name:  "Olahraga & Fitness",
			Type:  "expense",
			Color: stringPtr("#10B981"),
			Icon:  "Dumbbell",
		},
		{
			Name:  "Kesehatan",
			Type:  "expense",
			Color: stringPtr("#EF4444"),
			Icon:  "Heart",
		},
		{
			Name:  "Pendidikan",
			Type:  "expense",
			Color: stringPtr("#3B82F6"),
			Icon:  "BookOpen",
		},
		{
			Name:  "Lain-lain",
			Type:  "expense",
			Color: stringPtr("#6B7280"),
			Icon:  "MoreHorizontal",
		},
	}
	for _, category := range categories {
		if err := db.Where("name = ? AND type = ?", category.Name, category.Type).
			FirstOrCreate(&category).Error; err != nil {
			log.Printf("gagal membuat kategori %s: %v\n", category.Name, err)
		}
	}
}
