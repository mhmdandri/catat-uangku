package database

import (
	"catatan-keuangan/config"
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/attachments"
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/categories"
	groupmembers "catatan-keuangan/modules/group_members"
	"catatan-keuangan/modules/groups"
	"catatan-keuangan/modules/invitations"
	transactionlines "catatan-keuangan/modules/transaction_lines"
	"catatan-keuangan/modules/transactions"
	userprofile "catatan-keuangan/modules/user_profile"
	"catatan-keuangan/modules/users"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	c := config.Cfg
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal membuat koneksi ke database", err)
	}
	DB = db
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatal("Gagal mengaktifkan ekstensi uuid-ossp", err)
	}
	err = db.AutoMigrate(
		&users.User{},
		&accounts.Account{},
		&groups.Group{},
		&groupmembers.GroupMember{},
		&invitations.Invitation{},
		&categories.Category{},
		&transactions.Transactions{},
		&transactionlines.TransactionLine{},
		&attachments.Attachment{},
		&auth.RefreshToken{},
		&userprofile.UserProfile{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database", err)
	}
}
