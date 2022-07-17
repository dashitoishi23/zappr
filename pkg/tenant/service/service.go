package tenantservice

import (
	"errors"
	"time"

	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/constants"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gopkg.in/validator.v2"
)

type TenantService interface {
	AddTenant(newTenant tenantmodels.Tenant) (tenantmodels.Tenant, error)
}

type tenantService struct {
	repository repository.IRepository[tenantmodels.Tenant]
	roleRepository repository.IRepository[masterrolemodels.Role]
}

func NewService(repository repository.IRepository[tenantmodels.Tenant], 
roleRepository repository.IRepository[masterrolemodels.Role]) TenantService {
	return &tenantService{
		repository: repository,
		roleRepository: roleRepository,
	}
}

func (s *tenantService) AddTenant(newTenant tenantmodels.Tenant) (tenantmodels.Tenant, error) {
		if errs:= validator.Validate(newTenant); errs != nil {
			return newTenant, errors.New(constants.INVALID_MODEL)
		}

		tx := s.repository.GetTransaction()

		tx.Create(newTenant)

		normalUserRole := masterrolemodels.Role{
			Identifier: uuid.New().String(),
			Name: "Normal User",
			Scopes: pq.StringArray{
				"read",
			},
			TenantIdentifier: newTenant.Identifier,
			CreatedOn: time.Now(),
		}

		adminRole := masterrolemodels.Role{
			Identifier: uuid.New().String(),
			Name: "Admin",
			Scopes: pq.StringArray{
				"read",
				"write",
				"update",
				"delete",
			},
			TenantIdentifier: newTenant.Identifier,
			CreatedOn: time.Now(),
		}

		if err:= tx.Create(normalUserRole).Error; err!= nil {
			tx.Rollback()
			return newTenant, err
		}

		if err:= tx.Create(adminRole).Error; err != nil {
			tx.Rollback()
			return newTenant, err
		}

		if err:= tx.Commit().Error; err != nil {
			return newTenant, err
		}

		return newTenant, nil

}