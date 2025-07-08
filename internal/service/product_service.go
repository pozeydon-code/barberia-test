package service

import (
	"barberia-api-class/internal/domain"
	"context"
	"errors"
)

type ProductService struct {
	prodRepo domain.ProductRepo
}

func NewProductService(p domain.ProductRepo) *ProductService {
	return &ProductService{p}
}

func (s *ProductService) Create(ctx context.Context, prod *domain.Product) error {
	// 1. Validaciones
	if prod.Name == "" {
		return errors.New("el nombre del producto es requerido")
	}

	if prod.Price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}

	// 2. Verificar que no exista un producto con el mismo nombre
	existing, err := s.prodRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existingProd := range existing {
		if existingProd.Name == prod.Name {
			return errors.New("ya existe un producto con ese nombre")
		}
	}

	// 3. Crear
	return s.prodRepo.Create(ctx, prod)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	return s.prodRepo.GetById(ctx, id)
}

func (s *ProductService) ListAll(ctx context.Context) ([]domain.Product, error) {
	return s.prodRepo.List(ctx)
}

func (s *ProductService) Update(ctx context.Context, id uint, updatedProd *domain.Product) error {
	// 1. Validaciones
	if updatedProd.Name == "" {
		return errors.New("el nombre del producto es requerido")
	}

	if updatedProd.Price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}

	// 2. Verificar que el producto existe
	existing, err := s.prodRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// Verificar que no exista otro producto con el mismo nombre
	allProducts, err := s.prodRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existingProd := range allProducts {
		if existingProd.Name == updatedProd.Name && existingProd.ID != id {
			return errors.New("ya existe un producto con ese nombre")
		}
	}

	// Mantener ID original
	updatedProd.ID = existing.ID

	return s.prodRepo.Update(ctx, updatedProd)
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	// Verificar que el producto existe
	_, err := s.prodRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	return s.prodRepo.Delete(ctx, id)
}
