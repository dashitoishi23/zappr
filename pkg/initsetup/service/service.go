package initsetupservice

import (
	"encoding/json"
	"reflect"

	initsetupmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/models"
	repository "dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	util "dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
)

type InitSetupService interface {
	AddConfig(newConfig initsetupmodels.Config) (initsetupmodels.Config, error)
	EditConfig(newConfig initsetupmodels.Config) (bool, error)
}

type initSetupService struct {
	repository repository.IRepository[initsetupmodels.DBConfig]
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

	dbConfig.InitFields()
	tx := i.repository.Add(dbConfig)

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