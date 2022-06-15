package tenantservice

import (
	"errors"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/constants"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/repository"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type TenantService interface{
	CreateTenant(newTenant tenantmodels.Tenant) (tenantmodels.Tenant, error)
	GetTenant(identifier string) (tenantmodels.Tenant, error)
	UpdateTenant(updatedTenant tenantmodels.Tenant) (tenantmodels.Tenant, error)
	DeleteTenant(identifier string) bool
}

type tenantService struct{
	db *gorm.DB
	repository repository.IRepository[tenantmodels.Tenant]
}

func NewTenantService(database *gorm.DB) TenantService{
	return &tenantService{
		db: database,
		repository: repository.Repository[tenantmodels.Tenant](database),
	}
}

func (t *tenantService) CreateTenant(newTenant tenantmodels.Tenant) (tenantmodels.Tenant, error){
	if errs := validator.Validate(newTenant); errs != nil {
		return newTenant, errors.New(constants.INVALID_MODEL)
	}

	newTenant.Identifier = uuid.New().String()

	tx := t.repository.Add(newTenant)

	if tx.Error != nil {
		return newTenant, tx.Error
	}

	return newTenant, nil
}


