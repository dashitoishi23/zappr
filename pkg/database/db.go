package database

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func OpenDBConnection() (*gorm.DB, error) {
	user := os.Getenv("ZAPPR_POSTGRES_USER")
	password := os.Getenv("ZAPPR_POSTGRES_PASSWORD")
	database := os.Getenv("ZAPPR_POSTGRES_DB")
	host := os.Getenv("ZAPPR_POSTGRES_HOST")

	dsn := dsnBuilder(user, password, database, host, 5432)

	fmt.Print(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase: true,
		},
	})

	if err != nil {
		return db, err
	}

	return db, nil
}

func dsnBuilder(user string, password string, dbName string, host string, port int) string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, strconv.Itoa(port), dbName, user, password)
}
