package repositories

import (
	repositories "field-service/repositories/field"
	repositories2 "field-service/repositories/fieldschedule"
	repositories3 "field-service/repositories/time"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type RepositoryRegistryInterface interface {
	GetField() repositories.FieldRepositoriesInterface
	GetFieldSchedule() repositories2.FieldScheduleRepositoriesInterface
	GetTime() repositories3.TimeRepositoryInterface
}

func NewRegistryRepository(db *gorm.DB) RepositoryRegistryInterface {
	return &Registry{db: db}
}

func (r *Registry) GetField() repositories.FieldRepositoriesInterface {
	return repositories.NewFieldRepositories(r.db)
}

func (r *Registry) GetFieldSchedule() repositories2.FieldScheduleRepositoriesInterface {
	return repositories2.NewFieldScheduleRepositories(r.db)
}

func (r *Registry) GetTime() repositories3.TimeRepositoryInterface {
	return repositories3.NewTimeRepository(r.db)
}
