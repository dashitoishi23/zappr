package initsetupservice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	initsetupmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/models"
	repository "dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type InitSetupService interface {
	AddConfig(newConfig initsetupmodels.Config) (initsetupmodels.Config, error)
	EditConfig(newConfig initsetupmodels.Config) (bool, error)
}

type initSetupService struct {
	repository repository.IRepository[initsetupmodels.DBConfig]
	databaseConnection *gorm.DB
}

func NewInitSetupService(repository repository.IRepository[initsetupmodels.DBConfig]) InitSetupService {
	return &initSetupService{
		repository: repository,
	}
}

func (i *initSetupService) AddConfig(newConfig initsetupmodels.Config) (initsetupmodels.Config, error) {

	dbConfig, err := populateDBConfig(newConfig)

	if err != nil {
		return initsetupmodels.Config{}, err
	}

	errString:= i.isDBConnectionValid(newConfig)

	if errString != "Valid"{
		return initsetupmodels.Config{}, fmt.Errorf(errString)
	}

	errString = isRedisConnectionValid(newConfig)

	if errString != "Valid"{
		return initsetupmodels.Config{}, fmt.Errorf(errString)
	}

	fmt.Println("Database and Redis connections validated!")
	fmt.Println("Intialising DB")

	err = i.intialiseDB()

	if err != nil {
		return initsetupmodels.Config{}, err
	}

	dbConfig.InitFields()
	tx := i.repository.Add(dbConfig)

	if tx.Error != nil {
		return initsetupmodels.Config{}, tx.Error
	}

	db, _ := i.databaseConnection.DB()

	defer db.Close()

	return newConfig, tx.Error

}

func (i *initSetupService) EditConfig(newConfig initsetupmodels.Config) (bool, error) {
	dbConfig, err := populateDBConfig(newConfig)

	if err != nil {
		return false, err
	}

	existingConfig, err := i.repository.FindFirst(dbConfig)

	if err != nil {
		return false, err
	}

	dbConfig.UpdateFields(existingConfig.CreatedOn)
	
	_, err = i.repository.Update(dbConfig)

	return true, err
}

func populateDBConfig(newConfig initsetupmodels.Config) (initsetupmodels.DBConfig, error) {
	dbConfig := make(map[string]string)

	r := reflect.ValueOf(newConfig)

	for i:=0; i<r.NumField(); i++ {
		dbConfig[r.Field(i).String()] = r.Field(i).Interface().(string)
	}

	encodedJson, err := util.JsonEncoder(dbConfig)

	if err != nil {
		return initsetupmodels.DBConfig{}, err
	}

	return initsetupmodels.DBConfig{
		Config: json.RawMessage(encodedJson),
	}, nil

}

func (i *initSetupService) intialiseDB() (error) {
	initScript, err := os.ReadFile(filepath.Join("db_script.sql"))

	if err != nil {
		return err
	}

	tx := i.databaseConnection.Exec(string(initScript))

	if err != nil {
		return err
	}

	return tx.Error
}

func (i *initSetupService) isDBConnectionValid(newConfig initsetupmodels.Config) string {
	dsn := dsnBuilder(newConfig.DatabaseUser, newConfig.DatabasePassword, newConfig.DatabaseName, newConfig.DatabaseHost, newConfig.DatabasePort, newConfig.DatabaseSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	SkipDefaultTransaction: true,
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase: true,
		},
	})

	if err != nil{
		return err.Error()
	}

	database, err := db.DB()

	if err != nil {
		return err.Error()
	}

	err = database.Ping()

	if err != nil {
		return err.Error()
	}

	i.databaseConnection = db

	return "Valid"
}

func isRedisConnectionValid(newConfig initsetupmodels.Config) string {
	pool := &redis.Pool{
		MaxIdle: 80,
		MaxActive: 0,
		Dial: func() (redis.Conn, error) {

			if newConfig.RedisHost == "" {
				newConfig.RedisHost = "redis"
			}

			c, err := redis.Dial("tcp", newConfig.RedisHost + ":6379")
			
			if err != nil {
				panic(err)
			}

			if newConfig.RedisPassword != "" {
				_, authErr := c.Do("AUTH", newConfig.RedisPassword)

				if authErr != nil {
				panic(authErr)
			}
			}

			return c, err

		},
	}

	defer pool.Close()

	err := pool.Get().Err()

	if err != nil {
		return err.Error()
	}

	return "Valid"

}

func dsnBuilder(user string, password string, dbName string, host string, port int, sslMode string) string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, strconv.Itoa(port), dbName, user, password, sslMode)
}