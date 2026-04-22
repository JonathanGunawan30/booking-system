package controllers

import (
	"field-service/controllers/field"
	"field-service/controllers/fieldschedule"
	"field-service/controllers/time"
	"field-service/services"
)

type Registry struct {
	service services.ServiceRegistryInterface
}

type ControllerRegistryInterface interface {
	GetField() fieldcontroller.FieldControllerInterface
	GetFieldSchedule() fieldschedule.FieldScheduleControllerInterface
	GetTime() time.TimeControllerInterface
}

func NewControllerRegistry(service services.ServiceRegistryInterface) ControllerRegistryInterface {
	return &Registry{service: service}
}

func (r *Registry) GetField() fieldcontroller.FieldControllerInterface {
	return fieldcontroller.NewFieldController(r.service)
}

func (r *Registry) GetFieldSchedule() fieldschedule.FieldScheduleControllerInterface {
	return fieldschedule.NewFieldScheduleController(r.service)
}

func (r *Registry) GetTime() time.TimeControllerInterface {
	return time.NewTimeController(r.service)
}
