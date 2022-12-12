package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func OpenZapprConfigDBConnection() (*gorm.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("ZAPPR_POSTGRES_DB")
	host := os.Getenv("ZAPPR_POSTGRES_HOST")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")

	dsn := dsnBuilder(user, password, "postgres", host, 5432, sslMode)

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

	err = executeRawSQL(db, "db_create.sql")

	if err != nil {
		return db, err
	}

	dbVar, err := db.DB()

	if err != nil {
		return db, err
	}

	dbVar.Close()

	dsn = dsnBuilder(user, password, database, host, 5432, sslMode)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase: true,
		},
	})

	if err != nil {
		return db, err
	}

	err = executeRawSQL(db, "db_create_zapprdb.sql")

	if err != nil {
		return db, err
	}

	return db, nil
}

func dsnBuilder(user string, password string, dbName string, host string, port int, sslMode string) string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, strconv.Itoa(port), dbName, user, password, sslMode)
}

func executeRawSQL(db *gorm.DB, filename string) error {
	initScript, err := os.ReadFile(filepath.Join(filename))

	if err != nil {
		return err
	}

	tx := db.Exec(string(initScript))

	return tx.Error
}
