package service

import (
	"barberia-api-class/internal/domain"
	"context"
	"errors"
	"time"
)

type AppointmentService struct {
	apptRepo domain.AppointmentRepo
	prodRepo domain.ProductRepo
}

func NewAppointmentService(apptRepo domain.AppointmentRepo, prodRepo domain.ProductRepo) *AppointmentService {
	return &AppointmentService{
		apptRepo,
		prodRepo,
	}
}

func (s *AppointmentService) Schedule(ctx context.Context, appt *domain.Appointment) error {
	// 1. Validar que la cita sea en el futuro
	if appt.StartTime.Before(time.Now()) {
		return errors.New("la cita no puede ser en el pasado")
	}

	// 2. Validar que el tiempo de fin sea despues del tiempo de inicio
	if !appt.EndTime.After(appt.StartTime) {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// 3. Evitar solapamiento de turnos
	existingAppts, err := s.apptRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existingAppt := range existingAppts {
		if s.appointmentsOverlap(appt, &existingAppt) {
			return errors.New("la cita se solapa con otra existente")
		}
	}

	// 4. Validar que los productos existan
	for i, prod := range appt.Products {
		existingProd, err := s.prodRepo.GetById(ctx, prod.ID)
		if err != nil {
			if err == domain.ErrNotFound {
				return errors.New("producto no encontrado: " + prod.Name)
			}
			return err
		}

		// Reasignar el producto existente al slice de productos de la cita
		appt.Products[i] = *existingProd
	}

	// 5. Crear la cita
	return s.apptRepo.Create(ctx, appt)
}

func (s *AppointmentService) GetByID(ctx context.Context, id uint) (*domain.Appointment, error) {
	return s.apptRepo.GetById(ctx, id)
}

func (s *AppointmentService) ListAll(ctx context.Context) ([]domain.Appointment, error) {
	return s.apptRepo.List(ctx)
}

func (s *AppointmentService) Update(ctx context.Context, id uint, appt *domain.Appointment) error {
	// 1. Verificar que la cita exista
	existingAppt, err := s.apptRepo.GetById(ctx, appt.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			return errors.New("cita no encontrada")
		}
		return err
	}

	// 2. Mantener el ID Original
	appt.ID = existingAppt.ID

	// 3. Validar los horarios
	if appt.StartTime.Before(time.Now()) {
		return errors.New("la cita no puede ser en el pasado")
	}

	if !appt.EndTime.After(appt.StartTime) {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}

	// 4. Evitar solapamiento de turnos
	existingAppts, err := s.apptRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existing := range existingAppts {
		if existing.ID != id && s.appointmentsOverlap(appt, &existing) {
			return errors.New("la cita se solapa con otra existente")
		}
	}

	// Validar existencia de productos
	for i, prod := range appt.Products {
		existingProd, err := s.prodRepo.GetById(ctx, prod.ID)
		if err != nil {
			if err == domain.ErrNotFound {
				return errors.New("producto no encontrado: " + prod.Name)
			}
			return err
		}
		appt.Products[i] = *existingProd
	}

	return s.apptRepo.Update(ctx, appt)
}

func (s *AppointmentService) Cancel(ctx context.Context, id uint) error {
	// 1. Verificar que la cita exista
	existingAppt, err := s.apptRepo.GetById(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return errors.New("cita no encontrada")
		}
		return err
	}

	// 2. Validar que la cita no sea en el pasado
	if existingAppt.StartTime.Before(time.Now()) {
		return errors.New("no se puede cancelar una cita en el pasado")
	}

	// 3. Eliminar la cita
	return s.apptRepo.Delete(ctx, id)
}

func (s *AppointmentService) GetTotalPrice(ctx context.Context, id uint) (float64, error) {
	appt, err := s.apptRepo.GetById(ctx, id)
	if err != nil {
		return 0, err
	}

	var total float64

	for _, product := range appt.Products {
		total += product.Price
	}

	return total, nil
}

func (s *AppointmentService) appointmentsOverlap(appt1, appt2 *domain.Appointment) bool {
	return appt1.StartTime.Before(appt2.EndTime) && appt2.StartTime.Before(appt1.EndTime)
}
