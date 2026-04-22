package services

import (
	"field-service/common/cloudflare"
	"field-service/repositories"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
)

type Registry struct {
	repository repositories.RepositoryRegistryInterface
	r2         cloudflare.R2Client
}

type ServiceRegistryInterface interface {
	GetField() fieldService.FieldServiceInterface
	GetFieldSchedule() fieldScheduleService.FieldScheduleServiceInterface
	GetTime() timeService.TimeServiceInterface
}

func NewServiceRegistry(repository repositories.RepositoryRegistryInterface, r2 cloudflare.R2Client) ServiceRegistryInterface {
	return &Registry{r2: r2, repository: repository}
}

func (r *Registry) GetField() fieldService.FieldServiceInterface {
	return fieldService.NewFieldService(r.repository, r.r2)
}

func (r *Registry) GetFieldSchedule() fieldScheduleService.FieldScheduleServiceInterface {
	return fieldScheduleService.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() timeService.TimeServiceInterface {
	return timeService.NewTimeService(r.repository)
}
